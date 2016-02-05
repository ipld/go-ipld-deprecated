package coding

import (
	mc "github.com/jbenet/go-multicodec"
	mccbor "github.com/jbenet/go-multicodec/cbor"
	mcmux "github.com/jbenet/go-multicodec/mux"

	ipld "github.com/ipfs/go-ipld"
	pb "github.com/ipfs/go-ipld/coding/pb"
	memory "github.com/ipfs/go-ipld/memory"
)

// defaultCodec is the default applied if user does not specify a codec.
// Most new objects will never specify a codec. We track the codecs with
// the object so that multiple people using the same object will continue
// to marshal using the same codec. the only reason this is important is
// that the hashes must be the same.
var defaultCodec string

var muxCodec *mcmux.Multicodec

func init() {
	// by default, always encode things as cbor
	defaultCodec = string(mc.HeaderPath(mccbor.Header))
	muxCodec = mcmux.MuxMulticodec([]mc.Multicodec{
		CborMulticodec(),
		JsonMulticodec(),
		pb.Multicodec(),
	}, selectCodec)
}

// Multicodec returns a muxing codec that marshals to
// whatever codec makes sense depending on what information
// the IPLD object itself carries
func Multicodec() mc.Multicodec {
	return muxCodec
}

func selectCodec(v interface{}, codecs []mc.Multicodec) mc.Multicodec {
	vn, ok := v.(*memory.Node)
	if !ok {
		return nil
	}

	codecKey, err := codecKey(*vn)
	if err != nil {
		return nil
	}

	for _, c := range codecs {
		if codecKey == string(mc.HeaderPath(c.Header())) {
			return c
		}
	}

	return nil // no codec
}

func codecKey(n memory.Node) (string, error) {
	chdr, ok := (n)[ipld.CodecKey]
	if !ok {
		// if no codec is defined, use our default codec
		chdr = defaultCodec
		if pb.IsOldProtobufNode(n) {
			chdr = string(pb.Header)
		}
	}

	chdrs, ok := chdr.(string)
	if !ok {
		// if chdr is not a string, cannot read codec.
		return "", mc.ErrType
	}

	return chdrs, nil
}
