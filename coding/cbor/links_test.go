package cbor

import (
	"bytes"
	"testing"

	ipld "github.com/ipfs/go-ipld"
	reader "github.com/ipfs/go-ipld/coding/stream"
	readertest "github.com/ipfs/go-ipld/coding/stream/test"
	memory "github.com/ipfs/go-ipld/memory"
	multiaddr "github.com/jbenet/go-multiaddr"
	cbor "github.com/whyrusleeping/cbor/go"
)

func TestLinksStringEmptyMeta(t *testing.T) {
	var buf bytes.Buffer

	node, err := reader.NewNodeFromReader(memory.Node{
		ipld.LinkKey: "#/foo/bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = RawEncode(&buf, node.(ipld.NodeIterator), true)
	if err != nil {
		t.Error(err)
	}

	var expected []byte
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeTag, TagIPLDLink, nil)...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeText, uint64(len("#/foo/bar")), nil)...)
	expected = append(expected, []byte("#/foo/bar")...)

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Error("Incorrect encoding")
		t.Logf("Expected: %v", expected)
		t.Logf("Actual:   %v", buf.Bytes())
	}

	cbor, err := NewCBORDecoder(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Error(err)
	}

	readertest.CheckReader(t, cbor, []readertest.Callback{
		readertest.Cb(readertest.Path(), reader.TokenNode, nil),
		readertest.Cb(readertest.Path(), reader.TokenKey, ipld.LinkKey),
		readertest.Cb(readertest.Path(ipld.LinkKey), reader.TokenValue, "#/foo/bar"),
		readertest.Cb(readertest.Path(), reader.TokenEndNode, nil),
	})
}

func TestLinksStringNonEmptyMetaCheckOrdering(t *testing.T) {
	var buf bytes.Buffer

	node, err := reader.NewNodeFromReader(memory.Node{
		"size":       55,
		"00":         11,          // should be encoded first in the map (0 is before s and @ and is smaller)
		ipld.LinkKey: "#/foo/bar", // should come first, this is a link
	})
	if err != nil {
		t.Fatal(err)
	}

	err = RawEncode(&buf, node.(ipld.NodeIterator), true)
	if err != nil {
		t.Error(err)
	}

	var expected []byte
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeTag, TagIPLDLink, nil)...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeArray, 2, nil)...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeText, uint64(len("#/foo/bar")), nil)...)
	expected = append(expected, []byte("#/foo/bar")...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeMap, 2, nil)...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeText, uint64(len("00")), nil)...)
	expected = append(expected, []byte("00")...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeUint, 11, nil)...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeText, uint64(len("size")), nil)...)
	expected = append(expected, []byte("size")...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeUint, 55, nil)...)

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Error("Incorrect encoding")
		t.Logf("Expected: %v", expected)
		t.Logf("Actual:   %v", buf.Bytes())
	}

	cbor, err := NewCBORDecoder(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Error(err)
	}

	readertest.CheckReader(t, cbor, []readertest.Callback{
		readertest.Cb(readertest.Path(), reader.TokenNode, nil),
		readertest.Cb(readertest.Path(), reader.TokenKey, ipld.LinkKey),
		readertest.Cb(readertest.Path(ipld.LinkKey), reader.TokenValue, "#/foo/bar"),
		readertest.Cb(readertest.Path(), reader.TokenKey, "00"),
		readertest.Cb(readertest.Path("00"), reader.TokenValue, uint64(11)),
		readertest.Cb(readertest.Path(), reader.TokenKey, "size"),
		readertest.Cb(readertest.Path("size"), reader.TokenValue, uint64(55)),
		readertest.Cb(readertest.Path(), reader.TokenEndNode, nil),
	})
}

func TestLinksMultiAddr(t *testing.T) {
	var buf bytes.Buffer

	ma, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	if err != nil {
		t.Error(err)
		return
	}

	node, err := reader.NewNodeFromReader(memory.Node{
		ipld.LinkKey: ma.String(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = RawEncode(&buf, node.(ipld.NodeIterator), true)
	if err != nil {
		t.Error(err)
	}

	var expected []byte
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeTag, TagIPLDLink, nil)...)
	expected = append(expected, cbor.EncodeInt(cbor.MajorTypeBytes, uint64(len(ma.Bytes())), nil)...)
	expected = append(expected, ma.Bytes()...)

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Error("Incorrect encoding")
		t.Logf("Expected: %v", expected)
		t.Logf("Actual:   %v", buf.Bytes())
	}

	cbor, err := NewCBORDecoder(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Error(err)
	}

	readertest.CheckReader(t, cbor, []readertest.Callback{
		readertest.Cb(readertest.Path(), reader.TokenNode, nil),
		readertest.Cb(readertest.Path(), reader.TokenKey, ipld.LinkKey),
		readertest.Cb(readertest.Path(ipld.LinkKey), reader.TokenValue, ma.String()),
		readertest.Cb(readertest.Path(), reader.TokenEndNode, nil),
	})
}
