package sig

import (
	dag "github.com/ipfs/go-ipfsld/dag"
)

// Signature is an object that represents a cryptographic
// signature on any other merkle node.
//
// this would serialize to:
//
//   {
//    "@context": "/ipfs/<hash-of-schema>/signature"
//   	"key": "<hash1>",
//   	"object": "<hash2>",
//   	"sig": "<sign(sk, <hash2>)>"
//   }
//
type Signature struct {
	Key    dag.Link // the signing key
	Object dag.Link // what is signed
	Sig    []byte   // the data representing the signature
}

// Sign creates a signature from a given key and a link to data.
// Since this is a merkledag, signing the link is effectively the
// same as an hmac signature.
func Sign(key key.SigningKey, signed dag.Link) (Signature, error) {
	s := Signature{}
	s.Key = dag.LinkTo(key)
	s.Object = dag.LinkTo(signed)

	s, err := dag.Marshal(s.Object)
	if err != nil {
		return s, err
	}

	s.Sig, err = key.Sign(s)
	return s, err
}
