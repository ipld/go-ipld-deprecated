package ipldpb

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"reflect"
	"testing"

	mc "github.com/jbenet/go-multicodec"
	mcproto "github.com/jbenet/go-multicodec/protobuf"

	ipld "github.com/ipfs/go-ipld/memory"
)

var testfile []byte

func init() {
	var err error
	testfile, err = ioutil.ReadFile("testfile")
	if err != nil {
		panic("could not read testfile. please run: make testfile")
	}
}

func TestPBDecode(t *testing.T) {
	c := mcproto.Multicodec(&PBNode{})
	buf := bytes.NewBuffer(testfile)
	dec := c.Decoder(buf)

	// pass the /mdagv1
	if err := mc.ConsumeHeader(buf, Header); err != nil {
		t.Fatal("failed to consume header", err)
	}

	var pbn PBNode
	if err := dec.Decode(&pbn); err != nil {
		t.Fatal("failed to decode", err)
	}

	if len(pbn.Links) < 7 {
		t.Fatal("incorrect number of links")
	}
	if len(pbn.Data) == 0 {
		t.Error("should have some data")
	}

	findLink := func(s string) *PBLink {
		for _, l := range pbn.Links {
			if *l.Name == s {
				return l
			}
		}
		return nil
	}

	makefile := findLink("Makefile")
	if makefile == nil {
		t.Error("did not find Makefile")
	} else {
		if *makefile.Tsize < 700 || *makefile.Tsize > 4096 {
			t.Error("makefile incorrect size")
		}
	}
}

func TestPB2LD(t *testing.T) {
	buf := bytes.NewBuffer(testfile)
	dec := Multicodec().Decoder(buf)

	var n ipld.Node
	if err := dec.Decode(&n); err != nil {
		t.Fatal("failed to decode", err)
	}

	attrs, ok := n["@attrs"].(ipld.Node)
	if !ok {
		t.Log(n)
		t.Fatal("invalid ipld.@attrs")
	}

	data, ok := attrs["data"].([]byte)
	if !ok {
		t.Log(n)
		t.Fatal("invalid ipld.@attrs.data")
	}
	if len(data) == 0 {
		t.Error("should have some data")
	}

	links, ok := attrs["links"].([]ipld.Node)
	if !ok {
		t.Fatal("invalid ipld.@attrs.links")
	}
	if len(links) < 7 {
		t.Fatal("incorrect number of links")
	}

	findLink := func(s string) ipld.Node {
		for i, l := range links {
			s2, ok := l["name"].(string)
			if !ok {
				t.Log(l)
				t.Fatalf("invalid ipld.links[%d].name", i)
			}
			if s2 == s {
				return l
			}
		}
		return nil
	}

	makefileLink := findLink("Makefile")
	if makefileLink == nil {
		t.Error("did not find Makefile")
	}

	makefile := n["Makefile"].(ipld.Node)
	if makefile == nil {
		t.Error("did not find Makefile")
	} else {
		size, ok := makefile["size"].(uint64)
		if !ok {
			t.Log(makefile)
			t.Fatal("invalid ipld.links[makefile].size")
		}
		if size < 700 || size > 4096 {
			t.Log(makefile)
			t.Error("makefile incorrect size")
		}
		if !reflect.DeepEqual(makefile, makefileLink) {
			t.Error("makefile and @attrs.links[name=makefile] are not the same")
		}
	}
}

func TestLD2PB(t *testing.T) {
	decbuf := bytes.NewBuffer(testfile)
	encbuf := bytes.NewBuffer(nil)
	dec := Multicodec().Decoder(decbuf)
	enc := Multicodec().Encoder(encbuf)

	var n ipld.Node
	if err := dec.Decode(&n); err != nil {
		t.Fatal("failed to decode", err)
	}

	if err := enc.Encode(&n); err != nil {
		t.Log(n)
		t.Fatal("failed to encode", err)
	}

	if !bytes.Equal(testfile, encbuf.Bytes()) {
		t.Log(hex.Dump(testfile))
		t.Log(hex.Dump(encbuf.Bytes()))
		t.Fatal("decoded bytes != encoded bytes")
	}
}
