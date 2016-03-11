package json

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/ipfs/go-ipld/coding/base"
	reader "github.com/ipfs/go-ipld/coding/stream"
	mc "github.com/jbenet/go-multicodec"
)

const HeaderPath string = "/json"
const HeaderMsgioPath string = "/json/msgio"

var Header []byte
var HeaderMsgio []byte
var ErrAlreadyRead error = fmt.Errorf("Stream already read: unable to read it a second time")

func init() {
	Header = mc.Header([]byte(HeaderPath))
	HeaderMsgio = mc.Header([]byte(HeaderMsgioPath))
}

type JSONDecoder struct {
	r    io.Reader
	s    io.Seeker
	pos  int64
	lock sync.Mutex
}

type jsonParser struct {
	base.BaseReader
	decoder *json.Decoder
}

func NewJSONDecoder(r io.Reader) (*JSONDecoder, error) {
	var offset int64 = -1
	var err error

	s := r.(io.Seeker)
	if s != nil {
		offset, err = s.Seek(0, 1)
		if err != nil {
			return nil, err
		}
	}
	return &JSONDecoder{r, s, offset, sync.Mutex{}}, nil
}

func (d *JSONDecoder) Read(cb reader.ReadFun) error {
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
	jsonParser := &jsonParser{
		base.CreateBaseReader(cb),
		json.NewDecoder(d.r),
	}
	err := jsonParser.readValue()
	if err == reader.NodeReadAbort {
		err = nil
	}
	return err
}

func (p *jsonParser) readValue() error {
	token, err := p.decoder.Token()
	if err != nil {
		return err
	}
	//log.Printf("JSON: read token value %#v %T", token, token)
	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case '{':
			err = p.ExecCallback(reader.TokenNode, nil)
			if err != nil {
				p.Descope()
				return err
			}
			err = p.readNode()
			if err != nil {
				p.Descope()
				return err
			}
			err = p.ExecCallback(reader.TokenEndNode, nil)
			p.Descope()
			p.Descope()
			return err
			break
		case '[':
			err = p.ExecCallback(reader.TokenArray, nil)
			if err != nil {
				p.Descope()
				return err
			}
			err = p.readArray()
			if err != nil {
				p.Descope()
				return err
			}
			err = p.ExecCallback(reader.TokenEndArray, nil)
			p.Descope()
			p.Descope()
			return err
			break
		default:
			return fmt.Errorf("JSON: unexpected delimiter token %#v", token.(json.Delim))
		}
	} else {
		switch token.(type) {

		case json.Number:
			intValue, err1 := token.(json.Number).Int64()
			if err1 != nil {
				token = intValue
			} else {
				token, err = token.(json.Number).Float64()
				if err != nil {
					return fmt.Errorf("JSON: failed to convert %v to float64: %v", token.(json.Number), err)
				}
			}
		case float64:
			if sintValue := int(token.(float64)); token.(float64) == float64(sintValue) {
				token = sintValue
			} else if intValue := int64(token.(float64)); token.(float64) == float64(intValue) {
				token = intValue
			} else if uintValue := uint64(token.(float64)); token.(float64) == float64(uintValue) {
				token = uintValue
			}
		case string:
		case bool:
		case nil:
		default:
			return fmt.Errorf("JSON: Unexpected token %#v", token)
		}
		err := p.ExecCallback(reader.TokenValue, token)
		p.Descope()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *jsonParser) readNode() error {
	for {
		token, err := p.decoder.Token()
		if err != nil {
			return err
		}
		//log.Printf("JSON: read token node  %#v %T", token, token)

		if delim, ok := token.(json.Delim); ok && delim == '}' {
			return nil
		}

		strValue, isStr := token.(string)
		if !isStr {
			return fmt.Errorf("JSON: expect string for object key: got %#v", token)
		}
		err = p.ExecCallback(reader.TokenKey, strValue)
		if err != nil {
			p.Descope()
			return err
		}

		p.PushPath(strValue)
		err = p.readValue()
		p.PopPath()
		p.Descope()
		if err != nil {
			return err
		}
	}
}

func (p *jsonParser) readArray() error {
	var index uint64 = 0
	for {
		token, err := p.decoder.Token()
		if err != nil {
			return err
		}
		//log.Printf("JSON: read token array %#v %T", token, token)

		if delim, ok := token.(json.Delim); ok && delim == ']' {
			return nil
		}

		p.PushPath(index)
		err = p.readValue()
		p.PopPath()
		if err != nil {
			return err
		}

		index++
	}
}
