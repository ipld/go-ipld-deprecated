package memory

import (
	"testing"

	reader "github.com/ipfs/go-ipld/coding/stream"
	ldt "github.com/ipfs/go-ipld/coding/stream/test"
)

func TestReader(t *testing.T) {
	var node *Node

	node = &Node{
		"key":   "value",
		"items": []interface{}{"a", "b", "c"},
		"count": 3,
	}

	callbacks := []ldt.Callback{
		ldt.Cb(ldt.Path(), reader.TokenNode, nil),
		ldt.Cb(ldt.Path(), reader.TokenKey, "count"),
		ldt.Cb(ldt.Path("count"), reader.TokenValue, 3),
		ldt.Cb(ldt.Path(), reader.TokenKey, "items"),
		ldt.Cb(ldt.Path("items"), reader.TokenArray, nil),
		ldt.Cb(ldt.Path("items"), reader.TokenIndex, 0),
		ldt.Cb(ldt.Path("items", 0), reader.TokenValue, "a"),
		ldt.Cb(ldt.Path("items"), reader.TokenIndex, 1),
		ldt.Cb(ldt.Path("items", 1), reader.TokenValue, "b"),
		ldt.Cb(ldt.Path("items"), reader.TokenIndex, 2),
		ldt.Cb(ldt.Path("items", 2), reader.TokenValue, "c"),
		ldt.Cb(ldt.Path("items"), reader.TokenEndArray, nil),
		ldt.Cb(ldt.Path(), reader.TokenKey, "key"),
		ldt.Cb(ldt.Path("key"), reader.TokenValue, "value"),
		ldt.Cb(ldt.Path(), reader.TokenEndNode, nil),
	}

	ldt.CheckReader(t, node, callbacks)
}

func TestReaderSkip(t *testing.T) {
	var node *Node

	node = &Node{
		"key":   "value",
		"items": []interface{}{"a", "b", "c"},
		"count": 3,
	}

	callbacks := []ldt.Callback{
		ldt.Cb(ldt.Path(), reader.TokenNode, nil),
		ldt.Cb(ldt.Path(), reader.TokenKey, "count", reader.NodeReadSkip),
		ldt.Cb(ldt.Path(), reader.TokenKey, "items"),
		ldt.Cb(ldt.Path("items"), reader.TokenArray, nil),
		ldt.Cb(ldt.Path("items"), reader.TokenIndex, 0, reader.NodeReadSkip),
		ldt.Cb(ldt.Path("items"), reader.TokenIndex, 1),
		ldt.Cb(ldt.Path("items", 1), reader.TokenValue, "b"),
		ldt.Cb(ldt.Path("items"), reader.TokenIndex, 2),
		ldt.Cb(ldt.Path("items", 2), reader.TokenValue, "c"),
		ldt.Cb(ldt.Path("items"), reader.TokenEndArray, nil, reader.NodeReadSkip),
		ldt.Cb(ldt.Path(), reader.TokenKey, "key", reader.NodeReadAbort),
	}

	ldt.CheckReader(t, node, callbacks)
}
