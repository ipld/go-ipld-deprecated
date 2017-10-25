package cbor

import (
	"io"

	ipld "github.com/ipfs/go-ipld"
	ma "github.com/jbenet/go-multiaddr"
	cbor "github.com/whyrusleeping/cbor/go"
)

// Encode to CBOR, add the multicodec header
func Encode(w io.Writer, node ipld.NodeIterator, tags bool) error {
	_, err := w.Write(Header)
	if err != nil {
		return err
	}

	return RawEncode(w, node, tags)
}

func encodeFilter(val interface{}) interface{} {
	var ok bool
	var object map[string]interface{}
	var link interface{}
	var linkStr string
	var newObject map[string]interface{} = map[string]interface{}{}
	var linkPath interface{}
	var linkObject interface{}

	object, ok = val.(map[string]interface{})
	if !ok {
		return val
	}

	link, ok = object[ipld.LinkKey]
	if !ok {
		return val
	}

	linkStr, ok = link.(string)
	if !ok {
		return val
	}

	maddr, err := ma.NewMultiaddr(linkStr)
	for k, v := range object {
		if k != ipld.LinkKey {
			newObject[k] = v
		}
	}
	if err != nil || maddr.String() != linkStr {
		linkPath = linkStr
	} else {
		linkPath = maddr.Bytes()
	}

	if len(newObject) > 0 {
		linkObject = []interface{}{linkPath, newObject}
	} else {
		linkObject = linkPath
	}

	return &cbor.CBORTag{
		Tag:           TagIPLDLink,
		WrappedObject: linkObject,
	}
}

// Encode to CBOR, do not add the multicodec header
func RawEncode(w io.Writer, node ipld.NodeIterator, tags bool) error {
	enc := cbor.NewEncoder(w)
	if tags {
		enc.SetFilter(encodeFilter)
	}

	// Buffering to memory is absolutely necessary to normalize the CBOR data file
	// (pairs order in a map)
	mem, err := ipld.ToMemory(node)
	if err != nil {
		return err
	}

	return enc.Encode(mem)
}
