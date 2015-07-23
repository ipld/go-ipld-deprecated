package files

import (
	"bytes"
	"os"

	dag "github.com/ipfs/go-ipfsld/dag"
)

type StreamCombinator interface {
	Apply(ctx cxt.Context, root dag.Link, store dag.Store) (io.Reader, error)
}

// StackStreamCombinator is a combinator that executes
// stackstream code to produce binary output.
type StackStreamCombinator struct {
	// Code is a stackstream program
	Code []byte
}

func (ssc *StackStreamCombinator) Apply(ctx cxt.Context, root dag.Node, store dag.Store) (io.Reader, error) {
  ...
}

// StreamCombinator is an object that carries some code
// representing how to combine ipfs objects to produce
// It may carry:
// - Data: a raw data buffer
// - Chunks: links to other (sub)files
// - Combinator: function that produces output from Data and Chunks.
type File struct {
	Data       []byte
	Chunks     []dag.Link
	Combinator dag.Link // when in doubt, concat.
}
