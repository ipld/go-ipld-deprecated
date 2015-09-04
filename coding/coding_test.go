package ipfsld

import (
	"testing"

	ipld "github.com/ipfs/go-ipld"

	mctest "github.com/jbenet/go-multicodec/test"
)

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
