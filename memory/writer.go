package memory

import (
	"fmt"

	stream "github.com/ipfs/go-ipld/coding/stream"
)

type Writer struct {
	Node    Node
	stack   []interface{}
	curKey  string
	valPart []byte
}

func NewWriter() *Writer {
	return &Writer{Node{}, nil, "", nil}
}

func NewNodeFrom(r stream.NodeReader) (Node, error) {
	w := NewWriter()
	err := stream.Copy(r, w)
	return w.Node, err
}

func (w *Writer) WriteValue(val interface{}) error {
	if len(w.stack) == 0 {
		return fmt.Errorf("Cannot write value")
	}

	switch val.(type) {
	case string:
		val = string(append(w.valPart, []byte(val.(string))...))
	case []byte:
		val = append(w.valPart, val.([]byte)...)
	}

	last := len(w.stack) - 1
	switch w.stack[last].(type) {
	case *Node:
		(*w.stack[last].(*Node))[w.curKey] = val
	case *[]interface{}:
		*w.stack[last].(*[]interface{}) = append(*w.stack[last].(*[]interface{}), val)
	default:
		panic("Currupted stack")
	}

	w.curKey = ""
	w.valPart = nil
	return nil
}

func (w *Writer) WriteValuePart(val []byte) error {
	w.valPart = append(w.valPart, val...)
	return nil
}

func (w *Writer) WriteBeginNode(n_elems int) error {
	if len(w.stack) == 0 {
		w.stack = append(w.stack, &w.Node)
		return nil
	} else {
		n := Node{}
		err := w.WriteValue(n)
		w.stack = append(w.stack, &n)
		return err
	}
}

func (w *Writer) WriteNodeKey(key string) error {
	w.curKey = key
	return nil
}

func (w *Writer) WriteEndNode() error {
	if len(w.stack) == 0 {
		return fmt.Errorf("Cannot end node")
	}

	w.stack = w.stack[:len(w.stack)-1]
	return nil
}

func (w *Writer) WriteBeginArray(n_elems int) error {
	if len(w.stack) == 0 {
		return fmt.Errorf("Cannot start array")
	}

	slice := []interface{}{}
	w.stack = append(w.stack, w.curKey, &slice)
	return nil
}

func (w *Writer) WriteEndArray() error {
	if len(w.stack) <= 1 {
		return fmt.Errorf("Cannot end array")
	}

	slice := w.stack[len(w.stack)-1].(*[]interface{})
	w.curKey = w.stack[len(w.stack)-2].(string)
	w.stack = w.stack[:len(w.stack)-2]
	err := w.WriteValue(*slice)
	return err
}
