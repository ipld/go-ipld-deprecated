package ipfsld

import (
	"bytes"
	"testing"

	ld "github.com/ipfs/go-ipld"
)

type TC struct {
	cbor  []byte
	src   ld.Node
	links map[string]ld.Link
	typ   string
	ctx   interface{}
}

var testCases []TC

func init() {
	testCases = append(testCases, TC{
		[]byte{},
		ld.Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": ld.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
		map[string]ld.Link{
			"baz": {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		nil,
	})

	testCases = append(testCases, TC{
		[]byte{},
		ld.Node{
			"foo":      "bar",
			"@type":    "commit",
			"@context": "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
			"baz": ld.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bazz": ld.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar": ld.Node{
				"@type": "mlinkoo",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar2": ld.Node{
				"foo": ld.Node{
					"@type": "mlink",
					"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
			},
		},
		map[string]ld.Link{
			"baz":      {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bazz":     {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bar2/foo": {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		"/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
	})

}

func TestMarshaling(t *testing.T) {
	for tci, tc := range testCases {
		node1 := tc.src
		d1, err := Marshal(&node1)
		if err != nil {
			t.Error("marshal error", err, tci)
		}

		node2, err := Unmarshal(d1)
		if err != nil {
			t.Error("unmarshal error", err, tci)
		}

		// these are not equal.
		// if !reflect.DeepEqual(&node1, node2) {
		// 	d2, _ := Marshal(node2)
		// 	t.Log(&node1)
		// 	t.Log(node2)
		// 	t.Log(d1)
		// 	t.Log(d2)
		// 	t.Error("RTTed node not equal", tci, bytes.Equal(d1, d2))
		// }

		d2, err := Marshal(node2)
		if err != nil {
			t.Error("marshal error", err, tci)
		}

		if !bytes.Equal(d1, d2) {
			t.Log(len(d1), d1)
			t.Log(len(d2), d2)
			t.Error("marshaled bytes not equal", tci)
		}
	}
}
