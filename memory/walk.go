package memory

import (
	"errors"
	"path"
	"strconv"
	"strings"

	"github.com/ipfs/go-ipld/paths"
)

const pathSep = "/"

// SkipNode is a special value used with Walk and WalkFunc.
// If a WalkFunc returns SkipNode, the walk skips the curr
// node and its children. It behaves like file/filepath.SkipDir
var SkipNode = errors.New("skip node from Walk")

// WalkFunc is the type of the function called for each node
// visited by Walk. The root argument is the node from which
// the Walk began. The curr argument is the currently visited
// node. The path argument is the traversal path, from root
// to curr.
//
// If there was a problem walking to curr, the err argument
// will describe the problem and the function can decide
// how to handle the error (and Walk _will not_ descend into
// any of the children of curr).
//
// WalkFunc may return an error. If the error is the special
// SkipNode error, the children of curr are skipped. All other
// errors halt processing early. In this respect, it behaves
// just like file/filepath.WalkFunc
type WalkFunc func(root, curr Node, path string, err error) error

// Walk traverses the given root node and all its children, calling
// WalkFunc with every Node visited, including root. All errors
// that arise while visiting nodes are passed to given WalkFunc.
// The order in which children are visited is not deterministic.
// Walk traverses sequences as well, which is to mean the nodes
// below will be visted as "foo/0", "foo/1", and "foo/3":
//
//   { "foo": [
//     {"a":"aaa"}, // visited as foo/0
//     {"b":"bbb"}, // visited as foo/1
//     {"c":"ccc"}, // visited as foo/2
//   ]}
//
// Note Walk is purely local and does not traverse Links. For a
// version of Walk that does traverse links, see the ipld/traverse
// package.
func Walk(root Node, walkFn WalkFunc) error {
	return walk(root, root, "", walkFn)
}

// WalkFrom is just like Walk, but starts the Walk at given startFrom
// sub-node. It is the equivalent of a regular Walk call which skips
// all nodes which do not have startFrom as a prefix.
func WalkFrom(root Node, startFrom string, walkFn WalkFunc) error {
	start := GetPath(root, startFrom)
	if start == nil {
		return errors.New("no descendant at " + startFrom)
	}
	return walk(root, start, startFrom, walkFn)
}

// walk is used to implement Walk.
func walk(root Node, curr interface{}, npath string, walkFunc WalkFunc) error {

	if nc, ok := curr.(Node); ok { // it's a node!
		// first, call user's WalkFunc.
		err := walkFunc(root, nc, npath, nil)
		if err == SkipNode {
			return nil // ok, let's skip this one.
		} else if err != nil {
			return err // something bad happened, return early.
		}

		// then recurse.
		for k, v := range nc {
			// Skip empty path components
			if len(k) == 0 {
				continue
			}

			// skip any keys which contain "/" in them.
			// this is explicitly disallowed.
			if strings.Contains(k, pathSep) {
				continue
			}

			// skip keys starting with "@", it is reserved for directives
			// It can be escaped using "\@" in which case, "@" is not the first
			// character
			if k[0] == '@' {
				continue
			}

			k = paths.UnescapePathComponent(k)
			err := walk(root, v, path.Join(npath, k), walkFunc)
			if err != nil {
				return err
			}
		}

	} else if sc, ok := curr.([]interface{}); ok { // it's a slice!
		for i, v := range sc {
			k := strconv.Itoa(i)
			err := walk(root, v, path.Join(npath, k), walkFunc)
			if err != nil {
				return err
			}
		}

	} else { // it's just data.
		// ignore it.
	}
	return nil
}

// GetPath gets a descendant of root, at npath. GetPath
// uses the UNIX path abstraction: components of a
// path are delimited with "/". The path MUST start with "/".
func GetPath(root interface{}, path_ string) interface{} {
	path_ = path.Clean(path_)[1:] // skip root /
	return GetPathCmp(root, strings.Split(path_, pathSep))
}

// GetPathCmp gets a descendant of root, at npath.
func GetPathCmp(root interface{}, npath []string) interface{} {
	if len(npath) == 0 {
		return root // we're done.
	}
	if root == nil {
		return nil // nowhere to go
	}

	k := npath[0]
	if vn, ok := root.(Node); ok {
		// if node, recurse
		k = paths.EscapePathComponent(k)
		return GetPathCmp(vn[k], npath[1:])

	} else if vs, ok := root.([]interface{}); ok {
		// if slice, use key as an int offset
		i, err := strconv.Atoi(k)
		if err != nil {
			return nil
		}
		if i < 0 || i >= len(vs) { // nothing at such offset
			return nil
		}

		return GetPathCmp(vs[i], npath[1:])
	}

	return nil // cannot keep walking...
}
