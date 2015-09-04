package sig

import (
	"os"

	dag "github.com/ipfs/go-ipfsld/dag"
)

// Dir represents a directory in unixfs. The links are
// dag.Links with one additional property:
//  * unixMode - the full unix mode
//
// this would serialize to:
//
//    {
//      "@context": "/ipfs/<hash-of-schema>/unixdir",
//    	"<filename1>": { "@value": "<hash1>", "unixMode": <mode1> },
//    	"<filename2>": { "@value": "<hash2>", "unixMode": <mode2> },
//    }
//
type Dir map[string]dag.Link

func (d *Dir) Entry(e string) (dag.Link, error) {
	l, ok := d[e]
	if !ok {
		return nil, os.ErrNotExist
	}
	return l, nil
}

func (d *Dir) Mode(e string) (os.FileMode, error) {
	l, err := d.Entry(e)
	if err != nil {
		return 0, err
	}

	m, ok := l["unixMode"].(os.FileMode)
	if !ok {
		return 0, ErrInvalid
	}
	return m, nil
}
