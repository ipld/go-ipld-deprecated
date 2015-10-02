package jsonld

import (
	"testing"
	"reflect"

	ipld "github.com/ipfs/go-ipld"
)

type TC struct {
	src    ipld.Node
	jsonld ipld.Node
}

var testCases []TC

func init() {
	testCases = append(testCases, TC{
		src: ipld.Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": ipld.Node{
				"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
		jsonld: ipld.Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": ipld.Node{
				"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
	}, TC{
		src: ipld.Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"@container": "@index",
			"@index": "links",
			"baz": ipld.Node{
				"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
		jsonld: ipld.Node{
			"links": ipld.Node{
				"foo": "bar",
				"bar": []int{1, 2, 3},
				"baz": ipld.Node{
					"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
			},
		},
	}, TC{
		src: ipld.Node{
			"@attrs": ipld.Node{
				"attr": "val",
			},
			"foo":        "bar",
			"@index":     "files",
			"@type":      "commit",
			"@container": "@index",
			"@context": "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
			"baz": ipld.Node{
				"foobar": "barfoo",
				"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"\\@bazz": ipld.Node{
				"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar/ra\\b": ipld.Node{
				"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPb",
			},
			"bar": ipld.Node{
				"@container": "@index",
				"foo": ipld.Node{
					"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPa",
				},
			},
		},
		jsonld: ipld.Node{
			"attr": "val",
			"@type": "commit",
			"@context": "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
			"files": ipld.Node{
				"foo":        "bar",
				"baz": ipld.Node{
					"foobar": "barfoo",
					"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
				"@bazz": ipld.Node{
					"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
				"bar/ra\\b": ipld.Node{
					"mlink":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPb",
				},
				"bar": ipld.Node{
				},
			},
		},
	})
}

func TestParsing(t *testing.T) {
	for tci, tc := range testCases {
		t.Logf("===== Test case #%d =====", tci)
		doc := tc.src

		// check JSON-LD mode
		jsonld := ToLinkedDataAll(doc)
		if !reflect.DeepEqual(tc.jsonld, jsonld) {
			t.Errorf("JSON-LD version mismatch.\nGot:    %#v\nExpect: %#v", jsonld, tc.jsonld)
		} else {
			t.Log("JSON-LD version OK")
		}

	}
}
