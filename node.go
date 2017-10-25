package ipld

import (
	"errors"
	"math/big"
)

var errWrongType = errors.New("Incorrect type: could not convert")

// Represents a value that may be too large to read it all at once and that is
// split in multiple chunks.
type NodeValueIterator interface {
	// Read the next chunk and return it. Returns nil at the end.
	// Valid types for values are string and []byte
	// The function can return an error if the value vannot be read
	Value() (interface{}, error)

	// Abort reading the value
	Skip() error
}

// Represents a node object that is iterable.
// It iterates over key value pairs. When created, the NodeIterator does not
// have a current item and iteration must be started using Next()
type NodeIterator interface {
	// Go to next item. Sub items iteration is aborted automatically.
	// Returns true if there is an item
	// Returns false at the end of the iteration
	// In case of error, returns false and Error() returns the error
	Next() bool

	// Return true if the node is an object, false if it is an array. Always
	// available even before or after the iteration
	IsObject() bool

	// When Next() encounters an error, the error is available here.
	Error() error

	// Iterate the current item value
	// Return nil if the iteration hasn't started, is finished or if the current
	// item is not iterable
	Children() (NodeIterator, error)

	// Abort iteration. Does nothing if the iteration is finished. Also aborts
	// child iterators that were returned by the Children() method.
	Skip() error

	// Return the current item key, or nil if the iteration hasn't started or is
	// finished
	// This should return either a string or a uint64
	Key() interface{}

	// Return the current item value or nil if the iteration hasn't started, is
	// finished, or if the current item is iterable. In that case the Children()
	// method should be called instead.
	//
	// The types returned should only be:
	// - for numbers: int, int64, uint64, *"math/big".Int, float32, float64
	// - for strings: string, []byte
	// - for multi part strings: NodeValueIterator
	// - for booleans: bool
	// - for null: nil
	Value() (interface{}, error)
}

func ToNodeIterator(n interface{}) NodeIterator {
	return n.(NodeIterator)
}

func ToNodeValueIterator(n interface{}) NodeValueIterator {
	return n.(NodeValueIterator)
}

func ToMemory(n interface{}) (interface{}, error) {
	var v interface{}
	var err error

	switch n.(type) {
	case NodeIterator:
		i := n.(NodeIterator)

		if i.IsObject() {
			res := map[string]interface{}{}

			for i.Next() {
				v, err = i.Children()
				if err != nil {
					return res, err
				}

				if v == nil {
					v, err = i.Value()
					if err != nil {
						return res, err
					}
				}

				v, err = ToMemory(v)
				if err != nil {
					return res, err
				}

				res[i.Key().(string)] = v
			}

			return res, i.Error()

		} else {
			res := []interface{}{}

			for i.Next() {
				v, err = i.Children()
				if err != nil {
					return res, err
				}

				if v == nil {
					v, err = i.Value()
					if err != nil {
						return res, err
					}
				}

				v, err = ToMemory(v)
				if err != nil {
					return res, err
				}

				res = append(res, v)
			}

			return res, i.Error()
		}

	case NodeValueIterator:
		i := n.(NodeValueIterator)
		return ReadBytes(i)

	default:
		return n, nil
	}
}

func ReadBytes(i NodeValueIterator) ([]byte, error) {
	var v interface{}
	var err error
	var res []byte = []byte{}

	for v, err = i.Value(); v != nil; v, err = i.Value() {
		chunk := ToBytes(v)
		if chunk != nil {
			res = append(res, chunk...)
		}
	}

	return res, err
}

func ReadString(i NodeValueIterator) (string, error) {
	b, err := ReadBytes(i)
	return string(b), err
}

func ToBytes(n interface{}) []byte {
	var ok bool
	var res []byte
	switch n.(type) {
	case string:
		res = []byte(n.(string))
		ok = true
	case []byte:
		res = n.([]byte)
		ok = true
	default:
		ok = false
	}
	if res == nil {
		res = []byte{}
	}
	if ok {
		return res
	} else {
		return nil
	}
}

func ToBytesErr(n interface{}) ([]byte, error) {
	if it, ok := n.(NodeValueIterator); ok {
		return ReadBytes(it)
	} else if res := ToBytes(n); res != nil {
		return res, nil
	} else {
		return nil, errWrongType
	}
}

func ToString(n interface{}) *string {
	var ok bool
	var res string
	switch n.(type) {
	case string:
		res = n.(string)
		ok = true
	case []byte:
		res = string(n.([]byte))
		ok = true
	default:
		ok = false
	}
	if ok {
		return &res
	} else {
		return nil
	}
}

func ToStringErr(n interface{}) (string, error) {
	if it, ok := n.(NodeValueIterator); ok {
		return ReadString(it)
	} else if res := ToString(n); res != nil {
		return *res, nil
	} else {
		return "", errWrongType
	}
}

func ToUint(n interface{}) *uint64 {
	var res uint64
	var ok bool
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	switch n.(type) {
	case int:
		res = uint64(n.(int))
		ok = true
	case int64:
		res = uint64(n.(int64))
		ok = true
	case uint64:
		res = n.(uint64)
		ok = true
	case *big.Int:
		i := n.(*big.Int)
		if i.BitLen() <= 64 && i.Sign() >= 0 {
			res = i.Uint64()
			ok = true
		}
	default:
	}
	if ok {
		return &res
	} else {
		return nil
	}
}

func ToInt(n interface{}) *int64 {
	var res int64
	var ok bool
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	switch n.(type) {
	case int:
		res = int64(n.(int))
		ok = true
	case int64:
		res = n.(int64)
		ok = true
	case uint64:
		res = int64(n.(uint64))
		ok = true
	case *big.Int:
		i := n.(*big.Int)
		if i.BitLen() <= 63 {
			res = i.Int64()
			ok = true
		}
	default:
	}
	if ok {
		return &res
	} else {
		return nil
	}
}

func ToFloat(n interface{}) *float64 {
	var ok bool
	var res float64
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	switch n.(type) {
	case int:
		res = float64(n.(int))
		ok = true
	case int64:
		res = float64(n.(int64))
		ok = true
	case uint64:
		res = float64(n.(uint64))
		ok = true
	case float32:
		res = float64(n.(float32))
		ok = true
	case float64:
		res = n.(float64)
		ok = true
	case *big.Int:
		i := n.(*big.Int)
		if i.BitLen() <= 64 {
			if i.Sign() < 0 {
				res = -float64(big.NewInt(0).Abs(i).Uint64())
				ok = true
			} else {
				res = float64(i.Uint64())
				ok = true
			}
		}
	default:
	}
	if ok {
		return &res
	} else {
		return nil
	}
}
