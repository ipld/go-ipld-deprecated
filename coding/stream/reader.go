package stream

import (
	"errors"
	"fmt"
)

type nodeReadErrors struct{ error }

var NodeReadAbort error = &nodeReadErrors{errors.New("abort")}
var NodeReadSkip error = &nodeReadErrors{errors.New("skip")}

type ReaderToken int

const (
	Token0    ReaderToken = 0
	TokenNode ReaderToken = iota
	TokenKey
	TokenArray
	TokenIndex
	TokenValuePart
	TokenValue
	TokenEndNode
	TokenEndArray
)

func TokenName(tok ReaderToken) string {
	switch tok {
	case TokenNode:
		return "TokenNode"
	case TokenKey:
		return "TokenKey"
	case TokenArray:
		return "TokenArray"
	case TokenIndex:
		return "TokenIndex"
	case TokenValuePart:
		return "TokenValuePart"
	case TokenValue:
		return "TokenValue"
	case TokenEndNode:
		return "TokenEndNode"
	case TokenEndArray:
		return "TokenEndArray"
	default:
		return fmt.Sprintf("Token%d", tok)
	}
}

// Callback function to be called when reading a Node. The path represents the
// location in the hierarchy, tokenType is the type of token being read, and
// value is the corresponding value. The function can return SkipNode to skip
// recursing, or an error to stop the reading at any time.
type ReadFun func(path []interface{}, tokenType ReaderToken, value interface{}) error

// Represents a Node that can be read using the special Read function
type NodeReader interface {

	// Read the node from the beginning. Return any errors found during the
	// process. The function f is called for every value found in the node stream,
	// on the following events:
	//
	// - Once for each node when it starts (TokenNode) with a nil value
	// - Once for each key of each node (TokenKey) with the key name as value
	// - Once for each node when it closes (TokenEndNode) with a nil value
	// - Once for each array / slice when it starts (TokenArray) with a nil value
	// - Once for each slice item (TokenIndex) with the index as value
	// - Once for each array / slice when it stops (TokenEndArray) with a nil
	//   value
	// - Once for each other value (TokenValue) with the actual value given
	//
	// For example, the node that can be represented using the following JSON
	// construct:
	//
	//    {
	//      "key":   "value",
	//      "items": ["a", "b", "c"],
	//      "count": 3
	//    }
	//
	// Will result in the function called in this way:
	//
	//    f({},           TokenNode,     nil)
	//    f({},           TokenKey,      "key")
	//    f({"key"},      TokenValue,    "value")
	//    f({},           TokenKey,      "items")
	//    f({"items"},    TokenArray,    nil)
	//    f({"items"},    TokenIndex,    0)
	//    f({"items", 0}, TokenValue,    "a")
	//    f({"items"},    TokenIndex,    1)
	//    f({"items", 1}, TokenValue,    "b")
	//    f({"items"},    TokenIndex,    2)
	//    f({"items", 2}, TokenValue,    "c")
	//    f({"items"},    TokenEndArray, nil)
	//    f({},           TokenKey,      "count")
	//    f({"count"},    TokenValue,    3)
	//    f({},           TokenEndNode,  nil)
	//
	// Iteration order for maps is fixed by the order of the pairs in memory.
	//
	// Returning NodeReadSkip from the callback function will skip the sub elements:
	//
	// - for TokenNode or TokenArray, the complete node or array will be skipped
	// - for TokenKey or TokenIndex, the corresponding value will be skipped
	//
	// Returning NodeReadAbort will cancel the reading process without returning
	// an error to the caller.
	//
	// Returning any other error will abort the read process and the error will be
	// passed as a return value of the Read() function.
	//
	// When the value to read is too large, the value can be transmitted in
	// chunks. The chunks are all introduced by the TokenValuePart until the last
	// chunk which is introduced by TokenValue
	//
	// Allowed types for the path are:
	//
	// - string: for object keys
	// - int, uint64: for array indices
	//
	// Allowed value types are:
	//
	// - for numbers: int, int64, uint64, *"math/big".Int, float32, float64
	// - for strings: string, []byte
	// - for booleans: bool
	// - for null: nil
	//
	Read(f ReadFun) error
}
