package ipfsld

import (
	"bytes"
	"path"
	"strconv"

	mh "github.com/jbenet/go-multihash"
)

// Doc is an ipfs-ld document. effectively, it is just a JSON-LD document
// (which is {,de}serialized to JSON or CBOR) which derives from a base
// schema, the IPFS-LD schema (@context). This allows keys to specify:
//
//    "myfield": { "value": "Qmabcbcbdba", "@type": "mlink" }
//
// which is then taken to be a merkle-link, which IPFS handles specially.
type Doc struct {
	// Schema is a link to the schema of this document. If the document
	// does not define one, it is assumed to be the root IPFS-LD schema.
	Schema string

	// Data is the raw document data
	Data map[interface{}]interface{}

	// links is a map of string -> Link. this is a flattened map, using
	// JSON Pointer syntax. Meaning "foo/bar" points to "<hash>" in:
	//
	//   { "foo": {"bar": "<hash>" } }
	//
	links  map[string]Link
	linksP bool

	// sData is a map[string]interface{} for convenience
	sData  map[string]interface{}
	sDataP bool
}

func NewDoc(data map[interface{}]interface{}) *Doc {
	return &Doc{Schema: "placeholder", Data: data}
}

func (d *Doc) StrData() map[string]interface{} {
	if d.sDataP {
		return d.sData
	}

	for k, v := range d.Data {
		if sk, ok := k.(string); ok {
			d.sData[sk] = v
		}
	}
	d.sDataP = true
	return d.sData
}

func (d *Doc) Type() string {
	s, _ := d.StrData()["@type"].(string)
	return s
}

func (d *Doc) Context() interface{} {
	return d.StrData()["@context"]
}

// Links returns all the merkle-links in the document. When the document
// is parsed, all the links are identified and references are cached, so
// getting the links only walks the document _once_. Note though that the
// entire document must be walked.
func (d *Doc) Links() map[string]Link {
	if !d.linksP {
		d.links = Links(d)
		d.linksP = true
	}
	return d.links
}

// Link is a merkle-link structure.
type Link struct {
	Name string
	Hash mh.Multihash
	Data map[interface{}]interface{}
}

func (l Link) Equal(l2 Link) bool {
	if l.Name != l2.Name {
		return false
	}
	if !bytes.Equal(l.Hash, l2.Hash) {
		return false
	}
	// if !mapEqual(l.Data, l2.Data) {
	//   return false
	// }
	return true
}

func Links(doc *Doc) map[string]Link {
	m := map[string]Link{}
	walkDoc(m, "", doc.Data)
	return m
}

func walkDoc(links map[string]Link, jpath string, rest interface{}) {
	if rest == nil {
		return
	}

	if mrest, ok := rest.(map[interface{}]interface{}); ok { // it's a map!
		for k, v := range mrest {
			ks, ok := k.(string)
			if !ok {
				continue
			}

			jpath2 := path.Join(jpath, ks)
			if l, ok := getLink(jpath2, v); ok {
				links[jpath2] = l
			} else {
				walkDoc(links, jpath2, v)
			}
		}

	} else if arest, ok := rest.([]interface{}); ok { // it's an array!
		for i, v := range arest {
			jpath2 := path.Join(jpath, strconv.Itoa(i))
			if l, ok := getLink(jpath2, v); ok {
				links[jpath2] = l
			} else {
				walkDoc(links, jpath2, v)
			}
		}
	}
}

// checks whether a value is a link. for now we assume that all links
// follow:
//
//   { "linkname" : { "@type"  : "mlink", "hash": "<multihash>" } }
func isLink(value interface{}) bool {
	valmap, ok := value.(map[interface{}]interface{})
	if !ok {
		return false
	}

	ts, ok := valmap["@type"].(string)
	if !ok {
		return false
	}
	return ts == "mlink"
}

// returns the link value of an object. for now we assume that all links
// follow:
//
//   { "linkname" : { "@type"  : "mlink", "hash": "<multihash>" } }
func getLink(name string, value interface{}) (l Link, ok bool) {
	if !isLink(value) {
		return
	}

	valmap, ok := value.(map[interface{}]interface{})
	if !ok {
		return
	}
	hashS, ok := valmap["hash"].(string)
	if !ok {
		return
	}

	hash, err := mh.FromB58String(hashS)
	if err != nil {
		return
	}

	return Link{Name: name, Hash: hash, Data: valmap}, true
}
