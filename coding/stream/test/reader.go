package test

import (
	"testing"

	reader "github.com/ipfs/go-ipld/coding/stream"
	"github.com/mildred/assrt"
)

type Callback struct {
	Path      []interface{}
	TokenType reader.ReaderToken
	Value     interface{}
	Result    error
}

func Path(items ...interface{}) []interface{} {
	res := []interface{}{}
	for _, i := range items {
		res = append(res, i)
	}
	return res
}

func Cb(path []interface{}, tokenType reader.ReaderToken, value interface{}, res ...error) Callback {
	if len(res) == 0 {
		return Callback{path, tokenType, value, nil}
	} else if len(res) == 1 {
		return Callback{path, tokenType, value, res[0]}
	} else {
		panic("Cb() must be called with at most 4 arguments")
	}
}

func CheckReader(test *testing.T, node reader.NodeReader, callbacks []Callback) {
	assert := assrt.NewAssert(test)
	var i int = 0
	err := node.Read(func(path []interface{}, tokenType reader.ReaderToken, value interface{}) error {
		var e error
		if i >= len(callbacks) {
			assert.Logf("Callback %d: not described in test", i)
			assert.Logf("Should be: {%#v, %v, %#v}", path, reader.TokenName(tokenType), value)
			assert.Fail()
		} else {
			cb := callbacks[i]
			etk := reader.TokenName(cb.TokenType)
			atk := reader.TokenName(tokenType)
			actual := Callback{path, tokenType, value, cb.Result}
			assert.Logf("Callback %3d: %s %#v", i, etk, cb)
			assert.Logf("         got: %s %#v", atk, actual)
			assert.Equal(cb, actual)
			e = cb.Result
		}
		i++
		return e
	})
	assert.Equal(len(callbacks), i, "Number of callbacks incorrect")
	assert.Nil(err)
}
