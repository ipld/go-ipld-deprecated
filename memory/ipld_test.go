package memory

import (
	"testing"

	mh "github.com/jbenet/go-multihash"
)

type TC struct {
	src   Node
	links map[string]string
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
		src: Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": Node{
				"mlink": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"test": Node{
				// This is not a link because mlink is not a string but a Node
				"mlink": Node{
					"mlink": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
			},
		},
		links: map[string]string{
			"baz":        "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			"test/mlink": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
		},
		typ: "",
		ctx: nil,
	}, TC{
		src: Node{
			"@context": "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
			"baz": Node{
				"mlink": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bazz": Node{
				"mlink": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar": Node{
				"mlink": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPb",
			},
			"bar2": Node{
				"@bar": Node{
					"mlink": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPa",
				},
				"\\@foo": Node{
					"mlink": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPa",
				},
			},
		},
		links: map[string]string{
			"baz":       "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			"bazz":      "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			"bar":       "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPb",
			"bar2/@foo": "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPa",
		},
		typ: "",
		ctx: "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
	})
}

func TestParsing(t *testing.T) {
	for tci, tc := range testCases {
		t.Logf("===== Test case #%d =====", tci)
		doc := tc.src

		// check links
		links := doc.Links()
		t.Logf("links: %#v", links)
		if len(links) != len(tc.links) {
			t.Errorf("links do not match, not the same number of links, expected %d, got %d", len(tc.links), len(links))
		}
		for k, l1 := range tc.links {
			l2 := links[k]
			if l1 != l2["mlink"] {
				t.Errorf("links do not match. %d/%#v %#v != %#v[mlink]", tci, k, l1, l2)
			}
		}
	}
}
