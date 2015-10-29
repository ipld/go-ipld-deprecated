package ipfsld

import (
	"io/ioutil"
	"testing"
	"reflect"
	"bytes"

	ipld "github.com/ipfs/go-ipld"

	mc "github.com/jbenet/go-multicodec"
	mctest "github.com/jbenet/go-multicodec/test"
)

var json_testfile []byte

func init() {
	var err error
	json_testfile, err = ioutil.ReadFile("json.testfile")
	if err != nil {
		panic("could not read json.testfile. please run: make json.testfile")
	}
}

type TC struct {
	cbor  []byte
	src   ipld.Node
	links map[string]ipld.Link
	typ   string
	ctx   interface{}
}

var testCases []TC

func init() {
	testCases = append(testCases, TC{
		[]byte{},
		ipld.Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": ipld.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
		map[string]ipld.Link{
			"baz": {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		nil,
	})

	testCases = append(testCases, TC{
		[]byte{},
		ipld.Node{
			"foo":      "bar",
			"@type":    "commit",
			"@context": "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
			"baz": ipld.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bazz": ipld.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar": ipld.Node{
				"@type": "mlinkoo",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar2": ipld.Node{
				"foo": ipld.Node{
					"@type": "mlink",
					"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
			},
		},
		map[string]ipld.Link{
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
		var tcb ipld.Node
		mctest.RoundTripTest(t, codec, &(tca.src), &tcb)
	}
}

// Test decoding and encoding a json file
func TestJsonDecodeEncode(t *testing.T) {
	var n ipld.Node
	codec := Multicodec()

	if err := mc.Unmarshal(codec, json_testfile, &n); err != nil {
		t.Log(json_testfile)
		t.Error(err)
		return
	}

	linksExpected := map[string]ipld.Link{
		"abc": ipld.Link {
			"mlink": "QmXg9Pp2ytZ14xgmQjYEiHjVjMFXzCVVEcRTWJBmLgR39V",
		},
	}
	linksActual := ipld.Links(n)
	if !reflect.DeepEqual(linksExpected, linksActual) {
		t.Log(linksExpected)
		t.Log(linksActual)
		t.Logf("node: %#v\n", n)
		t.Fatalf("Links are not expected")
	}

	encoded, err := mc.Marshal(codec, &n)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(json_testfile, encoded) {
		t.Error("marshalled values not equal")
		t.Log(string(json_testfile))
		t.Log(string(encoded))
		t.Log(json_testfile)
		t.Log(encoded)
	}
}

