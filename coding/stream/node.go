package stream

import (
	ipld "github.com/ipfs/go-ipld"
	"sync/atomic"
)

type UnexpectedTokenError ReaderToken

func (e *UnexpectedTokenError) Error() string {
	return "Unexpected token " + TokenName(ReaderToken(*e))
}

type skipper interface {
	Skip() error
}

type args struct {
	path      []interface{}
	tokenType ReaderToken
	value     interface{}
	result    chan error
}

type valueIterator struct {
	c   chan *args
	e   *atomic.Value
	val *args
}

func (v *valueIterator) Value() (interface{}, error) {

	if v.c == nil {
		return nil, nil
	}

	if v.val == nil {
		v.val = <-v.c
		if v.val == nil {
			return nil, v.e.Load().(error)
		}
		close(v.val.result)
	}

	res := v.val.value
	v.val.result = nil

	switch v.val.tokenType {
	case TokenValue:
		v.val = nil
		v.c = nil
		return res, nil
	case TokenValuePart:
		v.val = nil
		return res, nil
	default:
		err := UnexpectedTokenError(v.val.tokenType)
		return nil, &err
	}
}

func (v *valueIterator) Skip() error {
	val, err := v.Value()
	for val != nil {
		val, err = v.Value()
	}
	return err
}

type iterator struct {
	c     chan *args
	e     *atomic.Value
	err   error
	key   *args
	val   *args
	child skipper
	kind  ReaderToken
}

func (i *iterator) Next() bool {
	if i.child != nil {
		i.err = i.child.Skip()
		i.child = nil
		if i.err != nil {
			return false
		}
	}

	if i.key != nil && i.key.result != nil {
		i.key.result <- NodeReadSkip
		close(i.key.result)
	}

	if i.c == nil {
		return false
	}

	i.key = <-i.c
	if i.key == nil {
		i.err = i.e.Load().(error)
		return false
	}

	i.val = nil
	switch i.key.tokenType {
	case TokenEndNode, TokenEndArray:
		close(i.key.result)
		i.c = nil
		i.key = nil
		i.val = nil
		return false
	case TokenKey, TokenIndex:
		break
	default:
		close(i.key.result)
		err := UnexpectedTokenError(i.key.tokenType)
		i.err = &err
		return false
	}

	return true
}

func (i *iterator) IsObject() bool {
	return i.kind == TokenNode
}

func (i *iterator) Error() error {
	return i.err
}

func (i *iterator) Children() (ipld.NodeIterator, error) {
	if i.child != nil || i.key == nil {
		// Already iterating, don't iterate twice
		// Or key not fetched yet
		// Or iteration end
		return nil, nil
	} else if i.key.result != nil {
		// Close key return channel
		close(i.key.result)
		i.key.result = nil
	}

	// Fetch the value of the key
	if i.val == nil {
		i.val = <-i.c
		if i.val == nil {
			return nil, i.e.Load().(error)
		}
		close(i.val.result)
		i.val.result = nil
	}

	switch i.val.tokenType {
	case TokenValuePart, TokenValue:
		return nil, nil
	case TokenArray, TokenNode:
		break
	default:
		err := UnexpectedTokenError(i.val.tokenType)
		return nil, &err
	}

	// create the child iterator
	res := new(iterator)
	res.c = i.c
	res.e = i.e
	res.kind = i.val.tokenType
	i.child = res
	return res, nil
}

func (i *iterator) Skip() error {
	for i.Next() {
		// Skip
	}
	return i.Error()
}

func (i *iterator) Key() interface{} {
	if i.key == nil {
		return nil
	} else if ii, ok := i.key.value.(int); ok && ii >= 0 {
		return uint64(ii)
	} else {
		return i.key.value
	}
}

func (i *iterator) Value() (interface{}, error) {
	if i.child != nil || i.key == nil {
		// Iterating over a child structure
		// Key not fetched yet or end of iteration
		return nil, nil
	} else if i.key.result != nil {
		// Close key return channel
		close(i.key.result)
		i.key.result = nil
	}

	// Fetch the value of the key
	if i.val == nil {
		i.val = <-i.c
		if i.val == nil {
			return nil, i.e.Load().(error)
		}
		close(i.val.result)
		i.val.result = nil
	}

	switch i.val.tokenType {
	case TokenValuePart:
		res := &valueIterator{i.c, i.e, i.val}
		i.child = res
		return res, nil
	case TokenValue:
		return i.val.value, nil
	case TokenArray, TokenNode:
		return nil, nil
	default:
		err := UnexpectedTokenError(i.val.tokenType)
		return nil, &err
	}
}

func (i *iterator) read(path []interface{}, tokenType ReaderToken, value interface{}) error {
	res := make(chan error)
	i.c <- &args{path, tokenType, value, res}
	return <-res
}

func NewNodeFromReader(r NodeReader) (interface{}, error) {
	var i *iterator = new(iterator)
	c := make(chan *args)
	i.c = c
	i.e = new(atomic.Value)
	go func() {
		err := r.Read(i.read)
		if err != nil {
			i.e.Store(err)
		}
		close(c)
	}()

	a := <-i.c
	if a == nil {
		return nil, i.e.Load().(error)
	}

	i.kind = a.tokenType

	close(a.result)
	switch a.tokenType {
	case TokenNode:
		return i, nil
	case TokenArray:
		return i, nil
	case TokenValue:
		return a.value, nil
	case TokenValuePart:
		return &valueIterator{i.c, i.e, a}, nil
	default:
		err := UnexpectedTokenError(a.tokenType)
		return nil, &err
	}
}
