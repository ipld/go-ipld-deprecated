package jsonld

import(
	ipld "github.com/ipfs/go-ipld"
)

const DefaultIndexName string = "@index"

func ContainerIndexName(n ipld.Node, defaultval string) string {
	var index_name string = defaultval

	index_val, ok := n["@index"]
	if str, is_string := index_val.(string); ok && is_string {
		index_name = str
	}

	return index_name
}

// Like ToLinkedDataAll but on the root node only, for use in Walk
func ToLinkedData(d ipld.Node) ipld.Node {
	attrs, directives, _, index := ipld.ParseNodeIndex(d)
	for k, v := range directives {
		if k != "@container" {
			attrs[k] = v
		}
	}
	if len(index) > 0 {
		index_name := ipld.ContainerIndexName(attrs, ipld.DefaultIndexName)
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

