package coding

import (
	"bytes"
	"fmt"
	"io"

	pb "github.com/ipfs/go-ipld/coding/pb"

	ipld "github.com/ipfs/go-ipld"
	stream "github.com/ipfs/go-ipld/coding/stream"
	mc "github.com/jbenet/go-multicodec"
)

var Header []byte

const (
	HeaderPath = "/mdagv1"
)

var StreamCodecs map[string]func(io.Reader) (stream.NodeReader, error)

type Codec int

const (
	NoCodec       Codec = 0
	CodecProtobuf Codec = iota
)

func init() {
	Header = mc.Header([]byte(HeaderPath))

	StreamCodecs = map[string]func(io.Reader) (stream.NodeReader, error){
		pb.MsgIOHeaderPath: func(r io.Reader) (stream.NodeReader, error) {
			return pb.Decode(mc.WrapHeaderReader(pb.MsgIOHeader, r))
		},
	}
}

func DecodeReader(r io.Reader) (stream.NodeReader, error) {
	// get multicodec first header, should be mcmux.Header
	err := mc.ConsumeHeader(r, Header)
	if err != nil {
		return nil, err
	}

	// get next header, to select codec
	hdr, err := mc.ReadHeader(r)
	if err != nil {
		return nil, err
	}

	hdrPath := string(mc.HeaderPath(hdr))

	fun, ok := StreamCodecs[hdrPath]
	if !ok {
		return nil, fmt.Errorf("no codec for %s", hdrPath)
	}
	return fun(r)
}

func Decode(r io.Reader) (interface{}, error) {
	rd, err := DecodeReader(r)
	if err != nil {
		return nil, err
	}

	return stream.NewNodeFromReader(rd)
}

func DecodeBytes(data []byte) (interface{}, error) {
	return Decode(bytes.NewReader(data))
}

func HasHeader(data []byte) bool {
	return len(data) >= len(Header) && bytes.Equal(data[:len(Header)], Header)
}

func DecodeLegacyProtobufBytes(data []byte) (stream.NodeReader, error) {
	return pb.RawDecode(data)
}

func EncodeRaw(codec Codec, w io.Writer, node ipld.NodeIterator) error {
	switch codec {
	case CodecProtobuf:
		return pb.Encode(w, node, true)
	default:
		return fmt.Errorf("Unknown codec %v", codec)
	}
}

func Encode(codec Codec, w io.Writer, node ipld.NodeIterator) error {
	w.Write(Header)
	return EncodeRaw(codec, w, node)
}

func EncodeBytes(codec Codec, node ipld.NodeIterator) ([]byte, error) {
	var buf bytes.Buffer
	err := Encode(codec, &buf, node)
	return buf.Bytes(), err
}
