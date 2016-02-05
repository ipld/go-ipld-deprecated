package memory

import (
	"errors"
	"path"
	"strconv"
)

// TransformFunc is the type of the function called for each node visited by
// Transform. The root argument is the node from which the Transform began. The
// curr argument is the currently visited node. The path argument is the
// traversal path, from root to curr.
//
// If there was a problem walking to curr, the err argument will describe the
// problem and the function can decide how to handle the error (and Transform
// _will not_ descend into any of the children of curr).
//
// TransformFunc may return a node, in which case the returned node will be used
// for further traversal instead of the curr node.
//
// TransformFunc may return an error. If the error is the special SkipNode
// error, the children of curr are skipped. All other errors halt processing
// early.
type TransformFunc func(root, curr Node, path []string, err error) (Node, error)

// Transform traverses the given root node and all its children, calling
// TransformFunc with every Node visited, including root. All errors that arise
// while visiting nodes are passed to given TransformFunc. The traversing
// algorithm is the same as the Walk function.
//
// Transform returns a node constructed from the different nodes returned by
// TransformFunc.
func Transform(root Node, transformFn TransformFunc) (Node, error) {
	n, err := transform(root, root, nil, transformFn)
	if node, ok := n.(Node); ok {
		return node, err
	} else {
		return nil, err
	}
}

// TransformFrom is just like Transform, but starts the Walk at given startFrom
// sub-node.
func TransformFrom(root Node, startFrom []string, transformFn TransformFunc) (interface{}, error) {
	start := GetPathCmp(root, startFrom)
	if start == nil {
		return nil, errors.New("no descendant at " + path.Join(startFrom...))
	}
	return transform(root, start, startFrom, transformFn)
}

// transform is used to implement Transform
func transform(root Node, curr interface{}, npath []string, transformFunc TransformFunc) (interface{}, error) {

	if nc, ok := curr.(Node); ok { // it's a node!
		// first, call user's WalkFunc.
		newnode, err := transformFunc(root, nc, npath, nil)
		res := Node{}
		if err == SkipNode {
			return newnode, nil // ok, let's skip this one.
		} else if err != nil {
			return nil, err // something bad happened, return early.
		} else if newnode != nil {
			nc = newnode
		}

		// then recurse.
		for k, v := range nc {
			n, err := transform(root, v, append(npath, k), transformFunc)
			if err != nil {
				return nil, err
			} else if n != nil {
				res[k] = n
			}
		}

		return res, nil

	} else if sc, ok := curr.([]interface{}); ok { // it's a slice!
		res := []interface{}{}
		for i, v := range sc {
			k := strconv.Itoa(i)
			n, err := transform(root, v, append(npath, k), transformFunc)
			if err != nil {
				return nil, err
			} else if n != nil {
				res = append(res, n)
			}
		}
		return res, nil

	} else { // it's just data.
		// ignore it.
	}
	return curr, nil
}
