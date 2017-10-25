package ipldpb

import (
	"fmt"
	"io"

	ipld "github.com/ipfs/go-ipld"
	memory "github.com/ipfs/go-ipld/memory"
	base58 "github.com/jbenet/go-base58"
	msgio "github.com/jbenet/go-msgio"
	mc "github.com/jbenet/go-multicodec"
)

const HeaderPath = "/mdagv1"
const MsgIOHeaderPath = "/protobuf/msgio"

var Header []byte
var MsgIOHeader []byte

var errInvalidLink = fmt.Errorf("invalid merkledag v1 protobuf, invalid Links")

func init() {
	Header = mc.Header([]byte(HeaderPath))
	MsgIOHeader = mc.Header([]byte(MsgIOHeaderPath))
}

func Decode(r io.Reader) (memory.Node, error) {
	err := mc.ConsumeHeader(r, MsgIOHeader)
	if err != nil {
		return nil, err
	}

	length, err := msgio.ReadLen(r, nil)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	return RawDecode(data)
}

func RawDecode(data []byte) (memory.Node, error) {
	var pbn *PBNode = new(PBNode)

	err := pbn.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	n := make(memory.Node)
	pb2ldNode(pbn, &n)

	return n, err
}

func Encode(w io.Writer, n ipld.NodeIterator, strict bool) error {
	_, err := w.Write(MsgIOHeader)
	if err != nil {
		return err
	}

	data, err := RawEncode(n, strict)
	if err != nil {
		return err
	}

	msgio.WriteLen(w, len(data))
	_, err = w.Write(data)
	return err
}

func RawEncode(n ipld.NodeIterator, strict bool) ([]byte, error) {
	pbn, err := ld2pbNode(n, strict)
	if err != nil {
		return nil, err
	}

	return pbn.Marshal()
}

func ld2pbNode(n ipld.NodeIterator, strict bool) (*PBNode, error) {
	var pbn PBNode
	has_data := false
	has_links := false

	for n.Next() {
		switch n.Key() {
		case "data":
			has_data = true
			data, err := n.Value()
			if err != nil {
				return nil, err
			}
			pbn.Data, err = ipld.ToBytesErr(data)
			if err != nil {
				return nil, err
			}
		case "links":
			has_links = true
			linkit, err := n.Children()
			if err != nil {
				return nil, err
			}
			for linkit.Next() {
				l, err := linkit.Children()
				if err != nil {
					return nil, err
				}
				pblink, err := ld2pbLink(l, strict)
				if err != nil {
					return nil, err
				}
				pbn.Links = append(pbn.Links, pblink)
			}
			if err := n.Error(); err != nil {
				return nil, err
			}
		default:
			if strict {
				return nil, fmt.Errorf("Invalid merkledag v1 protobuf: node contains extra field (%s)", n.Key())
			}
		}
	}

	if err := n.Error(); err != nil {
		return nil, err
	}

	if strict && !has_data {
		return nil, fmt.Errorf("Invalid merkledag v1 protobuf: no data")
	}

	if strict && !has_links {
		return nil, fmt.Errorf("Invalid merkledag v1 protobuf: no links")
	}

	return &pbn, nil
}

func pb2ldNode(pbn *PBNode, in *memory.Node) {
	var ordered_links []interface{}

	for _, link := range pbn.Links {
		ordered_links = append(ordered_links, pb2ldLink(link))
	}

	(*in)["data"] = pbn.GetData()
	(*in)["links"] = ordered_links
}

func pb2ldLink(pbl *PBLink) (link memory.Node) {
	defer func() {
		if recover() != nil {
			link = nil
		}
	}()

	link = make(memory.Node)
	link[ipld.LinkKey] = base58.Encode(pbl.Hash)
	link["name"] = *pbl.Name
	link["size"] = uint64(*pbl.Tsize)
	return link
}

func ld2pbLink(n ipld.NodeIterator, strict bool) (pbl *PBLink, err error) {
	pbl = &PBLink{}

	for n.Next() {
		switch n.Key() {
		case ipld.LinkKey:
			data, err := n.Value()
			if err != nil {
				return nil, err
			}
			hash := ipld.ToString(data)
			if hash == nil {
				return nil, fmt.Errorf("Invalid merkledag v1 protobuf: link is of incorect type")
			}
			pbl.Hash = base58.Decode(*hash)
			if strict && base58.Encode(pbl.Hash) != *hash {
				return nil, errInvalidLink
			}
		case "name":
			data, err := n.Value()
			if err != nil {
				return nil, err
			}
			name := ipld.ToString(data)
			if name == nil {
				return nil, fmt.Errorf("Invalid merkledag v1 protobuf: name is of incorect type")
			}
			pbl.Name = name
		case "size":
			data, err := n.Value()
			if err != nil {
				return nil, err
			}
			size := ipld.ToUint(data)
			if size == nil {
				return nil, fmt.Errorf("Invalid merkledag v1 protobuf: size is of incorect type")
			}
			pbl.Tsize = size
		default:
			if strict {
				return nil, fmt.Errorf("Invalid merkledag v1 protobuf: node contains extra field (%s)", n.Key())
			}
		}
	}

	return pbl, err
}
