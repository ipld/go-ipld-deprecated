package ipfsld

import (
	"bytes"
	"io"
	"reflect"

	codec "github.com/ugorji/go/codec"

	ld "github.com/ipfs/go-ipld"
)

var MapType reflect.Type

func init() {
	MapType = reflect.TypeOf(ld.Node(nil))
}

type Encoder interface {
	Encode(n *ld.Node) error
}

type Decoder interface {
	Decode(n *ld.Node) error
}

type encoder struct {
	enc *codec.Encoder
}

func NewEncoder(w io.Writer) *encoder {
	h := new(codec.CborHandle)
	h.MapType = MapType
	h.Canonical = true
	enc := codec.NewEncoder(w, h)
	return &encoder{enc}
}

func (c *encoder) Encode(n *ld.Node) error {
	return c.enc.Encode(&n)
}

type decoder struct {
	dec *codec.Decoder
}

func NewDecoder(r io.Reader) *decoder {
	h := new(codec.CborHandle)
	h.MapType = MapType
	h.Canonical = true
	dec := codec.NewDecoder(r, h)
	return &decoder{dec}
}

func (c *decoder) Decode(n *ld.Node) error {
	return c.dec.Decode(&n)
}

// Marshal serializes an ipfs-ld nument to a []byte.
func Marshal(n *ld.Node) ([]byte, error) {
	var buf bytes.Buffer
	err := MarshalTo(&buf, n)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalTo serializes an ipfs-ld nument to a writer.
func MarshalTo(w io.Writer, n *ld.Node) error {
	return NewEncoder(w).Encode(n)
}

// Unmarshal deserializes an ipfs-ld nument to a []byte.
func Unmarshal(buf []byte) (*ld.Node, error) {
	n := new(ld.Node)
	err := UnmarshalFrom(bytes.NewBuffer(buf), n)
	if err != nil {
		return nil, err
	}

	// have to call NewNode so the initial parsing (schema) takes place.
	return n, nil
}

// UnmarshalFrom deserializes an ipfs-ld nument from a reader.
func UnmarshalFrom(r io.Reader, n *ld.Node) error {
	return NewDecoder(r).Decode(n)
}
