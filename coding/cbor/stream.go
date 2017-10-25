package cbor

import (
	"fmt"
	"io"
	"math/big"
	"sync"

	ipld "github.com/ipfs/go-ipld"
	"github.com/ipfs/go-ipld/coding/base"
	reader "github.com/ipfs/go-ipld/coding/stream"
	ma "github.com/jbenet/go-multiaddr"
	mc "github.com/jbenet/go-multicodec"
	cbor "github.com/whyrusleeping/cbor/go"
)

const HeaderPath string = "/cbor"
const HeaderWithTagsPath string = "/cbor/ipld-tagsv1"

var Header []byte
var HeaderWithTags []byte
var ErrAlreadyRead error = fmt.Errorf("Stream already read: unable to read it a second time")

func init() {
	Header = mc.Header([]byte(HeaderPath))
	HeaderWithTags = mc.Header([]byte(HeaderWithTagsPath))
}

type CBORDecoder struct {
	r    io.Reader
	s    io.Seeker
	pos  int64
	lock sync.Mutex
}

type cborParser struct {
	base.BaseReader
	decoder  *cbor.Decoder
	tagStack []bool
}

func NewCBORDecoder(r io.Reader) (*CBORDecoder, error) {
	var offset int64 = -1
	var err error

	s := r.(io.Seeker)
	if s != nil {
		offset, err = s.Seek(0, 1)
		if err != nil {
			return nil, err
		}
	}
	return &CBORDecoder{r, s, offset, sync.Mutex{}}, nil
}

func (d *CBORDecoder) Read(cb reader.ReadFun) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.pos == -2 {
		return ErrAlreadyRead
	} else if d.pos == -1 {
		d.pos = -2
	} else {
		newoffset, err := d.s.Seek(d.pos, 0)
		if err != nil {
			return err
		} else if newoffset != d.pos {
			return fmt.Errorf("Failed to seek to position %d", d.pos)
		}
	}
	dec := cbor.NewDecoder(d.r)
	err := dec.DecodeAny(&cborParser{base.CreateBaseReader(cb), dec, []bool{false}})
	if err == reader.NodeReadAbort {
		err = nil
	}
	return err
}

func (p *cborParser) Prepare() error {
	return nil
}

func (p *cborParser) SetBytes(buf []byte) error {
	if p.tagStack[len(p.tagStack)-1] {
		maddr, err := ma.NewMultiaddrBytes(buf)
		if err != nil {
			return err
		}

		err = p.ExecCallback(reader.TokenKey, ipld.LinkKey)
		if err != nil {
			p.Descope()
			return err
		}
		p.PushPath(ipld.LinkKey)
		err = p.ExecCallback(reader.TokenValue, maddr.String())
		p.Descope()
		p.PopPath()
		p.Descope()
		return err
	}
	err := p.ExecCallback(reader.TokenValue, buf)
	p.Descope()
	return err
}

func (p *cborParser) SetUint(i uint64) error {
	err := p.ExecCallback(reader.TokenValue, i)
	p.Descope()
	return err
}

func (p *cborParser) SetInt(i int64) error {
	err := p.ExecCallback(reader.TokenValue, i)
	p.Descope()
	return err
}

func (p *cborParser) SetFloat32(f float32) error {
	err := p.ExecCallback(reader.TokenValue, f)
	p.Descope()
	return err
}

func (p *cborParser) SetFloat64(f float64) error {
	err := p.ExecCallback(reader.TokenValue, f)
	p.Descope()
	return err
}

func (p *cborParser) SetBignum(i *big.Int) error {
	err := p.ExecCallback(reader.TokenValue, i)
	p.Descope()
	return err
}

func (p *cborParser) SetNil() error {
	err := p.ExecCallback(reader.TokenValue, nil)
	p.Descope()
	return err
}

func (p *cborParser) SetBool(b bool) error {
	err := p.ExecCallback(reader.TokenValue, b)
	p.Descope()
	return err
}

func (p *cborParser) SetString(s string) error {
	if p.tagStack[len(p.tagStack)-1] {
		err := p.ExecCallback(reader.TokenKey, ipld.LinkKey)
		if err != nil {
			p.Descope()
			return err
		}
		p.PushPath(ipld.LinkKey)
		err = p.ExecCallback(reader.TokenValue, s)
		p.Descope()
		p.PopPath()
		p.Descope()
		return err
	}
	err := p.ExecCallback(reader.TokenValue, s)
	p.Descope()
	return err
}

func (p *cborParser) CreateMap() (cbor.DecodeValueMap, error) {
	p.tagStack = append(p.tagStack, false)
	if p.tagStack[len(p.tagStack)-2] {
		return p, nil
	}
	return p, p.ExecCallback(reader.TokenNode, nil)
}

func (p *cborParser) CreateMapKey() (cbor.DecodeValue, error) {
	return cbor.NewMemoryValue(""), nil
}

func (p *cborParser) CreateMapValue(key cbor.DecodeValue) (cbor.DecodeValue, error) {
	err := p.ExecCallback(reader.TokenKey, key.(*cbor.MemoryValue).Value)
	p.PushPath(key.(*cbor.MemoryValue).Value)
	return p, err
}

func (p *cborParser) SetMap(key, val cbor.DecodeValue) error {
	p.PopPath()
	p.Descope()
	return nil
}

func (p *cborParser) EndMap() error {
	p.tagStack = p.tagStack[:len(p.tagStack)-1]
	if p.tagStack[len(p.tagStack)-1] {
		return nil
	}

	err := p.ExecCallback(reader.TokenEndNode, nil)
	p.Descope()
	p.Descope()
	return err
}

func (p *cborParser) CreateArray(length int) (cbor.DecodeValueArray, error) {
	if p.tagStack[len(p.tagStack)-1] {
		return p, nil
	}
	return p, p.ExecCallback(reader.TokenArray, nil)
}

func (p *cborParser) GetArrayValue(index uint64) (cbor.DecodeValue, error) {
	if p.tagStack[len(p.tagStack)-1] {
		return p, nil
	}
	p.tagStack = append(p.tagStack, false)
	err := p.ExecCallback(reader.TokenIndex, index)
	p.PushPath(index)
	return p, err
}

func (p *cborParser) AppendArray(val cbor.DecodeValue) error {
	if p.tagStack[len(p.tagStack)-1] {
		return nil
	}
	p.PopPath()
	p.Descope()
	return nil
}

func (p *cborParser) EndArray() error {
	if p.tagStack[len(p.tagStack)-1] {
		return nil
	}
	err := p.ExecCallback(reader.TokenEndArray, nil)
	p.Descope()
	p.Descope()
	return err
}

func (p *cborParser) CreateTag(tag uint64, decoder cbor.TagDecoder) (cbor.DecodeValue, interface{}, error) {
	if tag == TagIPLDLink {
		p.tagStack = append(p.tagStack, true)
		return p, nil, p.ExecCallback(reader.TokenNode, nil)
	}
	return p, nil, nil
}

func (p *cborParser) SetTag(tag uint64, decval cbor.DecodeValue, decoder cbor.TagDecoder, val interface{}) error {
	if tag == TagIPLDLink {
		err := p.ExecCallback(reader.TokenEndNode, nil)
		p.Descope()
		p.Descope()
		p.tagStack = p.tagStack[:len(p.tagStack)-1]
		return err
	}
	return nil
}
