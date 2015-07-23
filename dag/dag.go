package merkledag

import (
	"fmt"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/golang.org/x/net/context"

	mh "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-multihash"
	key "github.com/ipfs/go-ipfs/blocks/key"
)

// Node represents a node in the IPFS Merkle DAG.
// nodes have opaque data and a set of navigable links.
type Node ld.Doc

// Links returns all the merkle-links in this node.
// This uses the ld.Links(...) function to find all
// then node's links.
func (n *Node) Links() map[string]Link {
	ld
}

// Link is a merkle-link to a target Node. The
// Link object is represented by a JSON-LD style
// map:
//
//   { "@type": "mlink", "@value": <multihash>, ... }
//
// Links support adding other data, which will be
// serialized and de-serialized along with the link.
// This allows users to set other properties on links:
//
//   {
//     "@type": "mlink",
//     "@value": <multihash>,
//     "unixType": "dir",
//     "unixMode": "0777",
//   }
//
// looking at a whole filesystem node, we might see something like:
//
//   {
//     "@context": "/ipfs/Qmf1ec6n9f8kW8JTLjqaZceJVpDpZD4L3aPoJFvssBE7Eb/merkleweb",
//     "foo": {
//       "@type": "mlink",
//       "@value": <multihash>,
//       "unixType": "dir",
//       "unixMode": "0777",
//     },
//     "bar": {
//       "@type": "mlink",
//       "@value": <multihash>,
//       "unixType": "file",
//       "unixMode": "0755",
//     }
//   }
//
type Link map[string]interface{}

// Value is the value of the link, the multihash
// it points to.
func (l *Link) Value() mh.Multihash {
	return l.Data["@value"]
}

// SetValue sets the multihash of the link.
func (l *Link) SetValue(h mh.Multihash) {
	l.Data["@value"] = h
}

// Type is the kind of link. For now, we only
// support "mlink"
func (l *Link) Type() string {
	return "mlink"
}

// Node retrieves a given node from a dagStore.
func (l *Link) Node(s *Store) (*Node, error) {
	return l.NodeCtx(context.Background(), s)
}

// Node retrieves a given node from a dagStore.
func (l *Link) NodeCtx(ctx cxt.Context, s *Store) (*Node, error) {
	return s.Get(ctx, Key(l.Value()))
}
