package stream

type NodeWriter interface {

	// Write a literal value. Accepted types are:
	// - int
	// - int64
	// - uint64
	// - float32
	// - float64
	// - *"math/big".Int
	// - string
	// - []byte
	WriteValue(val interface{}) error

	// Write a value part. When writing values in multiple chunks, one must call
	// WriteValuePart any number of times and end with WriteValue
	WriteValuePart(val []byte) error

	// Write the prolog for a node / associative array. n_elems is the number of
	// elements in the node, -1 if unknown.
	WriteBeginNode(n_elems int) error

	// Write a node key. Must be called before writing the value for that key.
	WriteNodeKey(key string) error

	// Write the epilog for a node
	WriteEndNode() error

	// Write the prolog for an array / slice. n_elems is the number of  elements
	// in the array, -1 if unknown.
	WriteBeginArray(n_elems int) error

	// Write the epilog for an array
	WriteEndArray() error
}
