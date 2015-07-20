package ipfsld

import (
	"reflect"
	"testing"

	mh "github.com/jbenet/go-multihash"
)

type TC struct {
	src   map[interface{}]interface{}
	links map[string]Link
	typ   string
	ctx   interface{}
}

var testCases []TC

func mmh(s string) mh.Multihash {
	h, err := mh.FromB58String(s)
	if err != nil {
		panic("invalid multihash: " + s)
	}
	return h
}

func init() {
	testCases = append(testCases, TC{
		map[interface{}]interface{}{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": map[interface{}]interface{}{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
		map[string]Link{
			"baz": {Name: "baz", Hash: mmh("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		nil,
	})

	testCases = append(testCases, TC{
		map[interface{}]interface{}{
			"foo":      "bar",
			"@type":    "commit",
			"@context": "/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
			"baz": map[interface{}]interface{}{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bazz": map[interface{}]interface{}{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar": map[interface{}]interface{}{
				"@type": "mlinkoo",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
			"bar2": map[interface{}]interface{}{
				"foo": map[interface{}]interface{}{
					"@type": "mlink",
					"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
				},
			},
		},
		map[string]Link{
			"baz":      {Name: "baz", Hash: mmh("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bazz":     {Name: "bazz", Hash: mmh("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bar2/foo": {Name: "bar2/foo", Hash: mmh("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		"/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
	})
}

func TestParsing(t *testing.T) {
	for tci, tc := range testCases {
		doc := NewDoc(tc.src)

		// check data
		if !reflect.DeepEqual(doc.Data, tc.src) {
			t.Errorf("data should be the same. %d", tci)
		}

		// check links
		links := doc.Links()
		for k, l1 := range tc.links {
			l2 := links[k]
			if !l1.Equal(l2) {
				t.Errorf("links do not match. %d/%s %s != %s", tci, k, l1, l2)
			}
		}
	}
}
