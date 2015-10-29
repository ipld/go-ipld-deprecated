package ipfsld

import (
	"io"

	mc "github.com/jbenet/go-multicodec"
	mcjson "github.com/jbenet/go-multicodec/json"
	ipld "github.com/ipfs/go-ipld"
)

type codec struct {
	mc.Multicodec
}

type decoder struct {
	mc.Decoder
}

func jsonMulticodec() mc.Multicodec {
	return &codec{mcjson.Multicodec(false)}
}

func (c *codec) Decoder(r io.Reader) mc.Decoder {
	return &decoder{ c.Multicodec.Decoder(r) }
}

func (c *decoder) Decode(v interface{}) error {
	err := c.Decoder.Decode(v)
	if err == nil {
		convert(v)
	}
	return err
}

func convert(val interface{}) interface{} {
	switch val.(type) {
	case *map[string]interface{}:
		vmi := val.(*map[string]interface{})
		n := ipld.Node{}
		for k, v := range *vmi {
			n[k] = convert(v)
			(*vmi)[k] = convert(v)
		}
		return &n
	case map[string]interface{}:
		vmi := val.(map[string]interface{})
		n := ipld.Node{}
		for k, v := range vmi {
			n[k] = convert(v)
			vmi[k] = convert(v)
		}
		return n
	case *[]interface{}:
		convert(*val.(*[]interface{}))
	case []interface{}:
		slice := val.([]interface{})
		for k, v := range slice {
			slice[k] = convert(v)
		}
	case *ipld.Node:
		convert(*val.(*ipld.Node))
	case ipld.Node:
		n := val.(ipld.Node)
		for k, v := range n {
			n[k] = convert(v)
		}
	default:
	}
	return val
}

