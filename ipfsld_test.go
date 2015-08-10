package ipld

import (
	"testing"

	mh "github.com/jbenet/go-multihash"
)

type TC struct {
	src   Node
	links map[string]Link
	typ   string
	ctx   interface{}
}

var testCases []TC

func mmh(b58 string) mh.Multihash {
	h, err := mh.FromB58String(b58)
	if err != nil {
		panic("failed to decode multihash")
	}
	return h
}

func init() {
	testCases = append(testCases, TC{
		Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
		map[string]Link{
			"baz": {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		nil,
	})

	testCases = append(testCases, TC{
		Node{
			"foo":      "bar",
			"@type":    "commit",
			"@context": "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
			"baz": Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bazz": Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar": Node{
				"@type": "mlinkoo",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar2": Node{
				"foo": Node{
					"@type": "mlink",
					"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
			},
		},
		map[string]Link{
			"baz":      {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bazz":     {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bar2/foo": {"@type": "mlink", "hash": ("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		"/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
	})
}

func TestParsing(t *testing.T) {
	for tci, tc := range testCases {
		doc := tc.src

		// check links
		links := doc.Links()
		t.Log(links)
		for k, l1 := range tc.links {
			l2 := links[k]
			if !l1.Equal(l2) {
				t.Errorf("links do not match. %d/%s %s != %s", tci, k, l1, l2)
			}
		}
	}
}
