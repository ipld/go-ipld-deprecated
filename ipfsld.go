package ipfsld

import (
	"errors"
	"path"
	"reflect"
	"strconv"
	"strings"

	mh "github.com/jbenet/go-multihash"
)

const (
	LinkType = "mlink"
	HashKey  = "hash"
	ValueKey = "@value"
	TypeKey  = "@type"
	CtxKey   = "@context"
	IDKey    = "@id"
)

// Doc is an ipfs-ld document. effectively, it is just a JSON-LD document
// (which is {,de}serialized to JSON or CBOR) which derives from a base
// schema, the IPFS-LD schema (@context). This allows keys to specify:
//
//    "myfield": { "value": "Qmabcbcbdba", "@type": "mlink" }
//
// which is then taken to be a merkle-link, which IPFS handles specially.
type Doc map[string]interface{}

func (d Doc) Get(path string) interface{} {
	return d.GetC(strings.Split(path, "/"))
}

func (d Doc) GetC(path []string) interface{} {
	if len(path) == 0 {
		return d
	}
	v := d[path[0]]
	if len(path) == 1 || v == nil {
		return d[path[0]]
	}
	if vd, ok := v.(Doc); ok {
		return vd.GetC(path[1:])
	}
	return v
}

func (d Doc) Type() string {
	s, _ := d[TypeKey].(string)
	return s
}

func (d Doc) Context() interface{} {
	return d[CtxKey]
}

// Links returns all the merkle-links in the document. When the document
// is parsed, all the links are identified and references are cached, so
// getting the links only walks the document _once_. Note though that the
// entire document must be walked.
func (d Doc) Links() map[string]Link {
	return Links(d)
}

// Link is a merkle-link structure.
type Link Doc

func (l Link) HashStr() string {
	s, _ := l[HashKey].(string)
	return s
}

func (l Link) Hash() (mh.Multihash, error) {
	s := l.HashStr()
	if s == "" {
		return nil, errors.New("no hash in link")
	}
	return mh.FromB58String(s)
}

func (l Link) Equal(l2 Link) bool {
	return reflect.DeepEqual(l, l2)
}

func Links(doc Doc) map[string]Link {
	m := map[string]Link{}
	walkDoc(m, "", doc)
	return m
}

func walkDoc(links map[string]Link, jpath string, rest interface{}) {
	if rest == nil {
		return
	}

	walkElem := func(k string, v interface{}) {
		jpath2 := path.Join(jpath, k)
		if l, ok := getLink(jpath2, v); ok {
			links[jpath2] = l
		} else {
			walkDoc(links, jpath2, v)
		}
	}

	if mrest, ok := rest.(Doc); ok { // it's a map!
		for k, v := range mrest {
			walkElem(k, v)
		}

	} else if mrest, ok := rest.(map[string]interface{}); ok { // it's a map!
		for k, v := range mrest {
			walkElem(k, v)
		}

	} else if arest, ok := rest.([]interface{}); ok { // it's an array!
		for i, v := range arest {
			walkElem(strconv.Itoa(i), v)
		}
	}
}

// checks whether a value is a link. for now we assume that all links
// follow:
//
//   { "linkname" : { "@type"  : "mlink", "hash": "<multihash>" } }
func isLink(value interface{}) bool {
	valmap, ok := value.(Doc)
	if !ok {
		return false
	}

	ts, ok := valmap[TypeKey].(string)
	if !ok {
		return false
	}
	return ts == LinkType
}

// returns the link value of an object. for now we assume that all links
// follow:
//
//   { "linkname" : { "@type"  : "mlink", "hash": "<multihash>" } }
func getLink(name string, value interface{}) (l Link, ok bool) {
	if !isLink(value) {
		return
	}

	l = make(Link)
	for k, v := range value.(Doc) {
		l[k] = v
	}
	return l, true
}
