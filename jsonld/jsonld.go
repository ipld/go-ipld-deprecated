package jsonld

import(
	ipld "github.com/ipfs/go-ipld"
)

const DefaultIndexName string = "@index"

func containerIndexName(n ipld.Node, defaultval string) string {
	var index_name string = defaultval

	index_val, ok := n["@index"]
	if str, is_string := index_val.(string); ok && is_string {
		index_name = str
	}

	return index_name
}

func isContainerIndex(n ipld.Node) bool {
	return n["@container"] == "@index"
}

// Like ToLinkedDataAll but on the root node only, for use in Walk
func ToLinkedData(d ipld.Node) ipld.Node {
	attrs, directives, _, index := ParseNodeIndex(d)
	for k, v := range directives {
		if k != "@container" {
			attrs[k] = v
		}
	}
	if len(index) > 0 {
		index_name := containerIndexName(attrs, DefaultIndexName)
		delete(attrs, "@index")
		if index_name[0] != '@' {
			attrs[index_name] = index
		}
	}
	return attrs
}

// Reorganize the data to be valid JSON-LD. This expand custom IPLD directives
// and unescape keys.
//
// The main processing now is to transform a IPLD data structure like this:
//
//	{
//		"@container": "@index",
//		"@index": "index-name",
//		"@attrs": {
//			"key": "value",
//		},
//		"index": { ... }
//	}
//
// to:
//
//	{
//		"key": "value",
//		"index-name": {
//			"index": { ... }
//		}
//	}
//
// In that case, it is good practice to define in the context the following
// type (this function cannot change the context):
//
//	"index-name": { "@container": "@index" }
//
func ToLinkedDataAll(d ipld.Node) ipld.Node {
	res, err := ipld.Transform(d, func(root, curr ipld.Node, path []string, err error) (ipld.Node, error) {
		return ToLinkedData(curr), err
	})
	if err != nil {
		panic(err) // should not happen
	}
	return res
}

func copyNode(n ipld.Node) ipld.Node {
	var res ipld.Node = ipld.Node{}
	for k, v := range n {
		res[k] = v
	}
	return res
}

func ParseNodeIndex(n ipld.Node) (attrs, directives, index ipld.Node, escapedIndex ipld.Node) {
	attrs = ipld.Node{}
	directives = ipld.Node{}
	index = ipld.Node{}
	escapedIndex = ipld.Node{}

	if real_attrs, ok := n["@attrs"]; ok {
		if attrs_node, ok := real_attrs.(ipld.Node); ok {
			attrs = copyNode(attrs_node)
		}
	}

	index_container := isContainerIndex(n)

	for key, val := range n {
		if key == "@attrs" {
			continue
		} else if key[0] == '@' {
			if key == "@index" {
				attrs[key] = val
			} else {
				directives[key] = val
			}
		} else {
			if index_container {
				escapedIndex[key] = val
				index[ipld.UnescapePathComponent(key)] = val
			} else {
				attrs[ipld.UnescapePathComponent(key)] = val
			}
		}
	}

	return
}

