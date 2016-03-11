package coding

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	ipld "github.com/ipfs/go-ipld"
	reader "github.com/ipfs/go-ipld/coding/stream"
	rt "github.com/ipfs/go-ipld/coding/stream/test"
	memory "github.com/ipfs/go-ipld/memory"
	assrt "github.com/mildred/assrt"
)

var codedFiles map[string][]byte = map[string][]byte{
	"json.testfile":     []byte{},
	"cbor.testfile":     []byte{},
	"protobuf.testfile": []byte{},
}

func init() {
	for fname := range codedFiles {
		var err error
		codedFiles[fname], err = ioutil.ReadFile(fname)
		if err != nil {
			panic("could not read " + fname + ". please run: make " + fname)
		}
	}
}

type TC struct {
	cbor  []byte
	src   memory.Node
	links map[string]memory.Link
	typ   string
	ctx   interface{}
}

// Test decoding and encoding a json and cbor file
func TestCodecsEncodeDecode(t *testing.T) {
	for fname, testfile := range codedFiles {

		r, err := DecodeBytes(testfile)
		if err != nil {
			t.Error(err)
			continue
		}

		var codec Codec
		switch fname {
		case "json.testfile":
			codec = CodecJSON
		case "cbor.testfile":
			codec = CodecCBOR
		case "protobuf.testfile":
			codec = CodecProtobuf
		default:
			panic("should not arrive here")
		}

		t.Logf("Decoded %s: %#v", fname, r)

		outData, err := EncodeBytes(codec, r.(ipld.NodeIterator))
		if err != nil {
			t.Error(err)
			continue
		}

		if !bytes.Equal(outData, testfile) {
			t.Errorf("%s: encoded is not the same as original", fname)
			t.Log(testfile)
			t.Log(string(testfile))
			t.Log(outData)
			t.Log(string(outData))
			f, err := os.Create(fname + ".error")
			if err != nil {
				t.Error(err)
			} else {
				defer f.Close()
				_, err := f.Write(outData)
				if err != nil {
					t.Error(err)
				}
			}
		}
	}
}

func TestJsonStream(t *testing.T) {
	a := assrt.NewAssert(t)
	t.Logf("Reading json.testfile")
	json, err := DecodeReader(bytes.NewReader(codedFiles["json.testfile"]))
	a.MustNil(err)

	rt.CheckReader(t, json, []rt.Callback{
		rt.Cb(rt.Path(), reader.TokenNode, nil),
		rt.Cb(rt.Path(), reader.TokenKey, "@codec"),
		rt.Cb(rt.Path("@codec"), reader.TokenValue, "/json"),
		rt.Cb(rt.Path(), reader.TokenKey, "abc"),
		rt.Cb(rt.Path("abc"), reader.TokenNode, nil),
		rt.Cb(rt.Path("abc"), reader.TokenKey, "mlink"),
		rt.Cb(rt.Path("abc", "mlink"), reader.TokenValue, "QmXg9Pp2ytZ14xgmQjYEiHjVjMFXzCVVEcRTWJBmLgR39V"),
		rt.Cb(rt.Path("abc"), reader.TokenEndNode, nil),
		rt.Cb(rt.Path(), reader.TokenEndNode, nil),
	})
}

func TestJsonStreamSkip(t *testing.T) {
	a := assrt.NewAssert(t)
	t.Logf("Reading json.testfile")
	json, err := DecodeReader(bytes.NewReader(codedFiles["json.testfile"]))
	a.MustNil(err)

	rt.CheckReader(t, json, []rt.Callback{
		rt.Cb(rt.Path(), reader.TokenNode, nil),
		rt.Cb(rt.Path(), reader.TokenKey, "@codec", reader.NodeReadSkip),
		rt.Cb(rt.Path(), reader.TokenKey, "abc"),
		rt.Cb(rt.Path("abc"), reader.TokenNode, nil),
		rt.Cb(rt.Path("abc"), reader.TokenKey, "mlink", reader.NodeReadAbort),
	})
}

func TestCborStream(t *testing.T) {
	a := assrt.NewAssert(t)
	t.Logf("Reading cbor.testfile")
	cbor, err := DecodeReader(bytes.NewReader(codedFiles["cbor.testfile"]))
	a.MustNil(err)

	rt.CheckReader(t, cbor, []rt.Callback{
		rt.Cb(rt.Path(), reader.TokenNode, nil),
		rt.Cb(rt.Path(), reader.TokenKey, "abc"),
		rt.Cb(rt.Path("abc"), reader.TokenNode, nil),
		rt.Cb(rt.Path("abc"), reader.TokenKey, "mlink"),
		rt.Cb(rt.Path("abc", "mlink"), reader.TokenValue, "QmXg9Pp2ytZ14xgmQjYEiHjVjMFXzCVVEcRTWJBmLgR39V"),
		rt.Cb(rt.Path("abc"), reader.TokenEndNode, nil),
		rt.Cb(rt.Path(), reader.TokenKey, "@codec"),
		rt.Cb(rt.Path("@codec"), reader.TokenValue, "/json"),
		rt.Cb(rt.Path(), reader.TokenEndNode, nil),
	})
}

func TestCborStreamSkip(t *testing.T) {
	a := assrt.NewAssert(t)
	t.Logf("Reading cbor.testfile")
	cbor, err := DecodeReader(bytes.NewReader(codedFiles["cbor.testfile"]))
	a.MustNil(err)

	rt.CheckReader(t, cbor, []rt.Callback{
		rt.Cb(rt.Path(), reader.TokenNode, nil),
		rt.Cb(rt.Path(), reader.TokenKey, "abc"),
		rt.Cb(rt.Path("abc"), reader.TokenNode, nil),
		rt.Cb(rt.Path("abc"), reader.TokenKey, "mlink", reader.NodeReadSkip),
		rt.Cb(rt.Path("abc"), reader.TokenEndNode, nil),
		rt.Cb(rt.Path(), reader.TokenKey, "@codec"),
		rt.Cb(rt.Path("@codec"), reader.TokenValue, "/json", reader.NodeReadAbort),
	})
}

func TestPbStream(t *testing.T) {
	a := assrt.NewAssert(t)
	t.Logf("Reading protobuf.testfile")
	t.Logf("Bytes: %v", codedFiles["protobuf.testfile"])
	pb, err := DecodeReader(bytes.NewReader(codedFiles["protobuf.testfile"]))
	a.MustNil(err)

	rt.CheckReader(t, pb, []rt.Callback{
		rt.Cb(rt.Path(), reader.TokenNode, nil),
		rt.Cb(rt.Path(), reader.TokenKey, "data"),
		rt.Cb(rt.Path("data"), reader.TokenValue, []byte{0x08, 0x01}),
		rt.Cb(rt.Path(), reader.TokenKey, "links"),
		rt.Cb(rt.Path("links"), reader.TokenArray, nil),
		rt.Cb(rt.Path("links"), reader.TokenIndex, 0),
		rt.Cb(rt.Path("links", 0), reader.TokenNode, nil),
		rt.Cb(rt.Path("links", 0), reader.TokenKey, ipld.LinkKey),
		rt.Cb(rt.Path("links", 0, ipld.LinkKey), reader.TokenValue, "Qmbvkmk9LFsGneteXk3G7YLqtLVME566ho6ibaQZZVHaC9"),
		rt.Cb(rt.Path("links", 0), reader.TokenKey, "name"),
		rt.Cb(rt.Path("links", 0, "name"), reader.TokenValue, "a"),
		rt.Cb(rt.Path("links", 0), reader.TokenKey, "size"),
		rt.Cb(rt.Path("links", 0, "size"), reader.TokenValue, uint64(10)),
		rt.Cb(rt.Path("links", 0), reader.TokenEndNode, nil),
		rt.Cb(rt.Path("links"), reader.TokenIndex, 1),
		rt.Cb(rt.Path("links", 1), reader.TokenNode, nil),
		rt.Cb(rt.Path("links", 1), reader.TokenKey, ipld.LinkKey),
		rt.Cb(rt.Path("links", 1, ipld.LinkKey), reader.TokenValue, "QmR9pC5uCF3UExca8RSrCVL8eKv7nHMpATzbEQkAHpXmVM"),
		rt.Cb(rt.Path("links", 1), reader.TokenKey, "name"),
		rt.Cb(rt.Path("links", 1, "name"), reader.TokenValue, "b"),
		rt.Cb(rt.Path("links", 1), reader.TokenKey, "size"),
		rt.Cb(rt.Path("links", 1, "size"), reader.TokenValue, uint64(10)),
		rt.Cb(rt.Path("links", 1), reader.TokenEndNode, nil),
		rt.Cb(rt.Path("links"), reader.TokenEndArray, nil),
		rt.Cb(rt.Path(), reader.TokenEndNode, nil),
	})
}

func TestPbStreamSkip(t *testing.T) {
	a := assrt.NewAssert(t)
	t.Logf("Reading protobuf.testfile")
	t.Logf("Bytes: %v", codedFiles["protobuf.testfile"])
	pb, err := DecodeReader(bytes.NewReader(codedFiles["protobuf.testfile"]))
	a.MustNil(err)

	rt.CheckReader(t, pb, []rt.Callback{
		rt.Cb(rt.Path(), reader.TokenNode, nil),
		rt.Cb(rt.Path(), reader.TokenKey, "data"),
		rt.Cb(rt.Path("data"), reader.TokenValue, []byte{0x08, 0x01}),
		rt.Cb(rt.Path(), reader.TokenKey, "links"),
		rt.Cb(rt.Path("links"), reader.TokenArray, nil),
		rt.Cb(rt.Path("links"), reader.TokenIndex, 0, reader.NodeReadSkip),
		rt.Cb(rt.Path("links"), reader.TokenIndex, 1),
		rt.Cb(rt.Path("links", 1), reader.TokenNode, nil),
		rt.Cb(rt.Path("links", 1), reader.TokenKey, ipld.LinkKey),
		rt.Cb(rt.Path("links", 1, ipld.LinkKey), reader.TokenValue, "QmR9pC5uCF3UExca8RSrCVL8eKv7nHMpATzbEQkAHpXmVM"),
		rt.Cb(rt.Path("links", 1), reader.TokenKey, "name", reader.NodeReadSkip),
		rt.Cb(rt.Path("links", 1), reader.TokenKey, "size"),
		rt.Cb(rt.Path("links", 1, "size"), reader.TokenValue, uint64(10), reader.NodeReadAbort),
	})
}
