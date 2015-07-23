package commit

import (
	dag "github.com/ipfs/go-ipfsld/dag"
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
	Parents   []dag.Link //
	Author    dag.Link   // link to an Authorship
	Committer dag.Link   // link to an Authorship
	Object    dag.Link   // what we version ("tree" in git)
	Comment   String     // describes the commit
}
