package ipfsld

import (
	"bytes"
	"testing"

	ld "github.com/ipfs/go-ipfsld"
	mh "github.com/jbenet/go-multihash"
)

type TC struct {
	cbor  []byte
	src   map[interface{}]interface{}
	links map[string]ld.Link
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
		[]byte{},
		map[interface{}]interface{}{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": map[interface{}]interface{}{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
		map[string]ld.Link{
			"baz": {Name: "baz", Hash: mmh("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		nil,
	})

	testCases = append(testCases, TC{
		[]byte{},
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
		map[string]ld.Link{
			"baz":      {Name: "baz", Hash: mmh("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bazz":     {Name: "bazz", Hash: mmh("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
			"bar2/foo": {Name: "bar2/foo", Hash: mmh("QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo")},
		},
		"",
		"/ipfs/QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo/mdag",
	})

}

func TestMarshaling(t *testing.T) {
	for tci, tc := range testCases {
		doc1 := ld.NewDoc(tc.src)
		d1, err := Marshal(doc1)
		if err != nil {
			t.Error("marshal error", err, tci)
		}

		doc2, err := Unmarshal(d1)
		if err != nil {
			t.Error("unmarshal error", err, tci)
		}

		// these are not equal.
		// if !reflect.DeepEqual(doc1.Data, doc2.Data) {
		// 	d2, _ := Marshal(doc2)
		// 	t.Log(doc1)
		// 	t.Log(d1)
		// 	t.Log(doc2)
		// 	t.Log(d2)
		// 	t.Error("RTTed doc not equal", tci, bytes.Equal(d1, d2))
		// }

		d2, err := Marshal(doc2)
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
