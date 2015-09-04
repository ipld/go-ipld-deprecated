package commit

import (
	ipld "github.com/ipfs/go-ipld"
)

// this would serialize to:
//
//   {
//     "@context": "/ipfs/<hash-to-commit-schema>/commit"
//     "parents": [ "<hash1>", ... ]
//     "author": "<hash2>",
//     "committer": "<hash3>",
//     "object": "<hash4>",
//     "comment": "comment as a string"
//   }
//
type Commit struct {
	Parents   []ipld.Link //
	Author    ipld.Link   // link to an Authorship
	Committer ipld.Link   // link to an Authorship
	Object    ipld.Link   // what we version ("tree" in git)
	Comment   String      // describes the commit
}

func (c *Commit) IPLDValidate() bool {
	// check at least one parent exists
	// check Parents have proper type
	// check author exists and has proper type
	// check commiter exists and has proper type
	// check object exists and has proper type
	return true
}
