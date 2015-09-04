package files

import (
	"bytes"
	"os"

	dag "github.com/ipfs/go-ipfsld/dag"
)

// File represents a readable byte stream.
// It may carry:
// - Data: a raw data buffer
// - Chunks: links to other (sub)files
// - Combinator: function that produces output from Data and Chunks.
type File struct {
	Data       []byte
	Chunks     []dag.Link
	Combinator dag.Link // when in doubt, concat.
}

// Reader returns an io.Reader which will read from
// the combinator.
func (f *File) Reader() (io.Reader, error) {
	// if no combinator is defined, output only Data.
	if f.Combinator == nil {
		return bytes.NewReader(f.Data)
	}

	l, err := d.Entry(e)
	if err != nil {
		return 0, err
	}

	m, ok := l["unixMode"]
	if !ok {
		return 0, ErrInvalid
	}

	mc, ok := m.(os.FileMode)
	if !ok {
		return 0, ErrInvalid
	}

	return m, nil
}

type File struct {
	Data   []byte
	Chunks []dag.Link
}

func Reader(f *File, s *dag.Store) (io.Reader, error) {
	return ReaderCtx(context.Background(), f, s)
}

func ReaderCtx(ctx cxt.Context, f *File, store *dag.Store) (io.Reader, error) {
	if f.Combinator == nil {
		return bytes.NewReader(f.Data), nil
	}

	var c Combinator
	if err := store.TypedGetCtx(ctx, f.Combinator.Hash(), &c); err != nil {
		return nil, err
	}

	dag.Unmarshal(cn)
	r := c.Combine()
}
