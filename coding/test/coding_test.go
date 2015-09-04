package ipfsld

import (
	"bytes"
	"testing"

	ld "github.com/ipfs/go-ipld"
	coding "github.com/ipfs/go-ipld/coding"
)

type TC struct {
	node ld.Node
}

var testCases []TC

func init() {
	testCases = append(testCases, TC{
		ld.Node{
			"foo": "bar",
			"bar": []int{1, 2, 3},
			"baz": ld.Node{
				"@type": "mlink",
				"hash":  "QmZku7P7KeeHAnwMr6c4HveYfMzmtVinNXzibkiNbfDbPo",
			},
		},
	})

	testCases = append(testCases, TC{
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
	})

}

func TMarshalRoundTrip(t *testing.T, c coding.Codec) {
	for tci, tc := range testCases {
		node1 := &tc.node
		d1, err := coding.Marshal(c, node1)
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

		d2, err := coding.Marshal(c, node2)
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

func TEncoderRoundTrip(t *testing.T, c coding.Codec) {
	var buf bytes.Buffer
	enc := c.Encoder(buf)
	dec := c.Decoder(buf)

	for tci, tc := range testCases {
		if err := enc.Encode(&tc.node); err != nil {
			t.Error("Encode error", tci, err)
		}
	}

	for tci, tc := range testCases {
		var node *ld.Node
		if err := enc.Decode(node); err != nil {
			t.Error("Decode error", tci, err)
		}

		if err := enc.

		d1, err := coding.Marshal(c, node1)
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

		d2, err := coding.Marshal(c, node2)
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
