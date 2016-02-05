package coding

import (
	"io"

	memory "github.com/ipfs/go-ipld/memory"
	mc "github.com/jbenet/go-multicodec"
	mccbor "github.com/jbenet/go-multicodec/cbor"
	mcjson "github.com/jbenet/go-multicodec/json"
)

type transformCodec struct {
	mc.Multicodec
}

type transformDecoder struct {
	mc.Decoder
}

func JsonMulticodec() mc.Multicodec {
	return &transformCodec{mcjson.Multicodec(false)}
}

func CborMulticodec() mc.Multicodec {
	return &transformCodec{mccbor.Multicodec()}
}

func (c *transformCodec) Decoder(r io.Reader) mc.Decoder {
	return &transformDecoder{c.Multicodec.Decoder(r)}
}

func (c *transformDecoder) Decode(v interface{}) error {
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
		n := memory.Node{}
		for k, v := range *vmi {
			n[k] = convert(v)
			(*vmi)[k] = convert(v)
		}
		return &n
	case map[string]interface{}:
		vmi := val.(map[string]interface{})
		n := memory.Node{}
		for k, v := range vmi {
			n[k] = convert(v)
			vmi[k] = convert(v)
		}
		return n
	case *map[interface{}]interface{}:
		vmi := val.(*map[interface{}]interface{})
		n := memory.Node{}
		for k, v := range *vmi {
			if k2, ok := k.(string); ok {
				n[k2] = convert(v)
				(*vmi)[k2] = convert(v)
			}
		}
		return &n
	case map[interface{}]interface{}:
		vmi := val.(map[interface{}]interface{})
		n := memory.Node{}
		for k, v := range vmi {
			if k2, ok := k.(string); ok {
				n[k2] = convert(v)
				vmi[k2] = convert(v)
			}
		}
		return n
	case *[]interface{}:
		convert(*val.(*[]interface{}))
	case []interface{}:
		slice := val.([]interface{})
		for k, v := range slice {
			slice[k] = convert(v)
		}
	case *memory.Node:
		convert(*val.(*memory.Node))
	case memory.Node:
		n := val.(memory.Node)
		for k, v := range n {
			n[k] = convert(v)
		}
	default:
	}
	return val
}
