package stream

import (
	"fmt"
)

func Copy(r NodeReader, w NodeWriter) error {
	return r.Read(func(path []interface{}, tk ReaderToken, value interface{}) error {
		switch tk {
		case TokenNode:
			return w.WriteBeginNode(-1)
		case TokenKey:
			return w.WriteNodeKey(value.(string))
		case TokenArray:
			return w.WriteBeginArray(-1)
		case TokenIndex:
			return nil
		case TokenValuePart:
			switch value.(type) {
			case string:
				return w.WriteValuePart([]byte(value.(string)))
			case []byte:
				return w.WriteValuePart(value.([]byte))
			default:
				return fmt.Errorf("Cannot write value part %#v (wrong type)", value)
			}
		case TokenValue:
			return w.WriteValue(value)
		case TokenEndNode:
			return w.WriteEndNode()
		case TokenEndArray:
			return w.WriteEndArray()
		default:
			return fmt.Errorf("Unexpected token %s", TokenName(tk))
		}
	})
}
