package ipldpb

import (
	"errors"
	"fmt"
	"io"

	mc "github.com/jbenet/go-multicodec"
	mcproto "github.com/jbenet/go-multicodec/protobuf"

	ipld "github.com/ipfs/go-ipld"
)

var Header []byte

var (
	errInvalidData = fmt.Errorf("invalid merkledag v1 protobuf, Data not bytes")
	errInvalidLink = fmt.Errorf("invalid merkledag v1 protobuf, invalid Links")
)

func init() {
	Header = mc.Header([]byte("/mdagv1"))
}

type codec struct {
	pbc mc.Multicodec
}

func Multicodec() mc.Multicodec {
	var n *PBNode
	return &codec{mcproto.Multicodec(n)}
}

func (c *codec) Encoder(w io.Writer) mc.Encoder {
	return &encoder{w: w, c: c, pbe: c.pbc.Encoder(w)}
}

func (c *codec) Decoder(r io.Reader) mc.Decoder {
	return &decoder{r: r, c: c, pbd: c.pbc.Decoder(r)}
}

func (c *codec) Header() []byte {
	return Header
}

type encoder struct {
	w   io.Writer
	c   *codec
	pbe mc.Encoder
}

type decoder struct {
	r   io.Reader
	c   *codec
	pbd mc.Decoder
}

func (c *encoder) Encode(v interface{}) error {
	nv, ok := v.(*ipld.Node)
	if !ok {
		return errors.New("must encode *ipld.Node")
	}

	if _, err := c.w.Write(c.c.Header()); err != nil {
		return err
	}

	n, err := ld2pbNode(nv)
	if err != nil {
		return err
	}

	return c.pbe.Encode(n)
}

func (c *decoder) Decode(v interface{}) error {
	nv, ok := v.(*ipld.Node)
	if !ok {
		return errors.New("must decode to *ipld.Node")
	}

	if err := mc.ConsumeHeader(c.r, c.c.Header()); err != nil {
		return err
	}

	var pbn PBNode
	if err := c.pbd.Decode(&pbn); err != nil {
		return err
	}

	pb2ldNode(&pbn, nv)
	return nil
}

func ld2pbNode(in *ipld.Node) (*PBNode, error) {
	n := *in
	var pbn PBNode
	var attrs ipld.Node

	if attrsvalue, hasattrs := n["@attrs"]; hasattrs {
		var ok bool
		attrs, ok = attrsvalue.(ipld.Node)
		if !ok {
			return nil, errInvalidData
		}
	} else {
		return &pbn, nil
	}

	if data, hasdata := attrs["data"]; hasdata {
		data, ok := data.([]byte)
		if !ok {
			return nil, errInvalidData
		}
		pbn.Data = data
	}

	if links, haslinks := attrs["links"]; haslinks {
		links, ok := links.([]ipld.Node)
		if !ok {
			return nil, errInvalidLink
		}

		for _, link := range links {
			pblink := ld2pbLink(link)
			if pblink == nil {
				return nil, fmt.Errorf("%s (%s)", errInvalidLink, link["name"])
			}
			pbn.Links = append(pbn.Links, pblink)
		}
	}
	return &pbn, nil
}

func pb2ldNode(pbn *PBNode, in *ipld.Node) {
	*in = make(ipld.Node)
	n := *in

	links := make([]ipld.Node, len(pbn.Links))
	for i, link := range pbn.Links {
		links[i] = pb2ldLink(link)
		n[ipld.EscapePathComponent(link.GetName())] = links[i]
	}

	n["@attrs"] = ipld.Node{
		"links": links,
		"data": pbn.Data,
	}
}

func pb2ldLink(pbl *PBLink) (link ipld.Node) {
	defer func() {
		if recover() != nil {
			link = nil
		}
	}()

	link = make(ipld.Node)
	link["hash"] = pbl.Hash
	link["name"] = *pbl.Name
	link["size"] = uint64(*pbl.Tsize)
	return link
}

func ld2pbLink(link ipld.Node) (pbl *PBLink) {
	defer func() {
		if recover() != nil {
			pbl = nil
		}
	}()

	hash := link["hash"].([]byte)
	name := link["name"].(string)
	size := link["size"].(uint64)

	pbl = &PBLink{}
	pbl.Hash = hash
	pbl.Name = &name
	pbl.Tsize = &size
	return pbl
}

func IsOldProtobufNode(n ipld.Node) bool {
	if len(n) > 2 { // short circuit
		return false
	}

	links, hasLinks := n["links"]
	_, hasData := n["data"]

	switch len(n) {
	case 2: // must be links and data
		if !hasLinks || !hasData {
			return false
		}
	case 1: // must be links or data
		if !(hasLinks || hasData) {
			return false
		}
	default: // nope.
		return false
	}

	if len(n) > 2 {
		return false // only links and data.
	}

	if hasLinks {
		links, ok := links.([]ipld.Node)
		if !ok {
			return false // invalid links.
		}

		// every link must be a mlink
		for _, link := range links {
			if !ipld.IsLink(link) {
				return false
			}
		}
	}

	return true // ok looks like an old protobuf node
}
