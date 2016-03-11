package json

import (
	"encoding/json"
	"io"

	ipld "github.com/ipfs/go-ipld"
)

// Encode to JSON, add the multicodec header
func Encode(w io.Writer, node ipld.NodeIterator) error {
	_, err := w.Write(Header)
	if err != nil {
		return err
	}

	return RawEncode(w, node)
}

// Encode to JSON, do not add the multicodec header
func RawEncode(w io.Writer, node ipld.NodeIterator) error {
	mem, err := ipld.ToMemory(node)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(mem)
}
