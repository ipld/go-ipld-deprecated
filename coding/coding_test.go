package coding

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	memory "github.com/ipfs/go-ipld/memory"

	mc "github.com/jbenet/go-multicodec"
	mctest "github.com/jbenet/go-multicodec/test"
)

var codedFiles map[string][]byte = map[string][]byte{
	"json.testfile": []byte{},
	"cbor.testfile": []byte{},
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

var testCases []TC

func init() {
	testCases = append(testCases, TC{
		[]byte{},
		memory.Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": memory.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
		map[string]memory.Link{
			"baz": {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		nil,
	})

	testCases = append(testCases, TC{
		[]byte{},
		memory.Node{
			"foo":      "bar",
			"@type":    "commit",
			"@context": "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
			"baz": memory.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bazz": memory.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar": memory.Node{
				"@type": "mlinkoo",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar2": memory.Node{
				"foo": memory.Node{
					"@type": "mlink",
					"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
			},
		},
		map[string]memory.Link{
			"baz":      {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bazz":     {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bar2/foo": {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		"/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
	})

}

func TestHeaderMC(t *testing.T) {
	codec := Multicodec()
	for _, tc := range testCases {
		mctest.HeaderTest(t, codec, &tc.src)
	}
}

func TestRoundtripBasicMC(t *testing.T) {
	codec := Multicodec()
	for _, tca := range testCases {
		var tcb memory.Node
		mctest.RoundTripTest(t, codec, &(tca.src), &tcb)
	}
}

// Test decoding and encoding a json and cbor file
func TestCodecsDecodeEncode(t *testing.T) {
	for fname, testfile := range codedFiles {
		var n memory.Node
		codec := Multicodec()

		if err := mc.Unmarshal(codec, testfile, &n); err != nil {
			t.Log(testfile)
			t.Error(err)
			continue
		}

		linksExpected := map[string]memory.Link{
			"abc": memory.Link{
				"mlink": "QmXg9Pp2ytZ14xgmQjYEiHjVjMFXzCVVEcRTWJBmLgR39V",
			},
		}
		linksActual := memory.Links(n)
		if !reflect.DeepEqual(linksExpected, linksActual) {
			t.Logf("Expected: %#v", linksExpected)
			t.Logf("Actual:   %#v", linksActual)
			t.Logf("node: %#v\n", n)
			t.Error("Links are not expected in " + fname)
			continue
		}

		encoded, err := mc.Marshal(codec, &n)
		if err != nil {
			t.Error(err)
			return
		}

		if !bytes.Equal(testfile, encoded) {
			t.Error("marshalled values not equal in " + fname)
			t.Log(string(testfile))
			t.Log(string(encoded))
			t.Log(testfile)
			t.Log(encoded)
		}
	}
}
