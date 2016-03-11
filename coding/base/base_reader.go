package base

import (
	"github.com/ipfs/go-ipld/coding/stream"
)

type BaseReader struct {
	Callback  stream.ReadFun
	cbEnabled []bool
	path      []interface{}
	tokens    []stream.ReaderToken
}

func CreateBaseReader(cb stream.ReadFun) BaseReader {
	return BaseReader{cb, []bool{}, []interface{}{}, []stream.ReaderToken{}}
}

// Executes the callback and stay in scope for sub elements
func (p *BaseReader) ExecCallback(token stream.ReaderToken, value interface{}) error {
	var err error
	enabled := !p.Skipping()
	if enabled {
		err = p.Callback(p.path, token, value)
		if err == stream.NodeReadSkip {
			enabled = false
			err = nil
		}
	}
	p.cbEnabled = append(p.cbEnabled, enabled)
	return err
}

// Return true if a parent callback wants to skip processing of its children
func (p *BaseReader) Skipping() bool {
	enabled := true
	if len(p.cbEnabled) > 0 {
		enabled = p.cbEnabled[len(p.cbEnabled)-1]
	}
	return !enabled
}

// Must be called after all sub elements below a ExecCallback are processed
func (p *BaseReader) Descope() {
	p.cbEnabled = p.cbEnabled[:len(p.cbEnabled)-1]
}

// Push a path element
func (p *BaseReader) PushPath(elem interface{}) {
	p.path = append(p.path, elem)
}

// Pop a path element
func (p *BaseReader) PopPath() {
	p.path = p.path[:len(p.path)-1]
}
