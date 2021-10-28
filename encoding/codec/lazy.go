package codec

import (
	"errors"
)

var (
	ErrLazyWrongValueType = errors.New("codec lazy converts on wrong type value")
)

type Lazy struct{}

func init() {
	Register(new(Lazy))
}

func (*Lazy) String() string {
	return "lazy"
}

func (*Lazy) Marshal(v interface{}) (Bytes, error) {
	switch v := v.(type) {
	case Bytes:
		return v.Dulplicate(), nil
	case *Bytes:
		return v.Dulplicate(), nil
	default:
		return nilBytes, ErrLazyWrongValueType
	}
}

func (s *Lazy) Unmarshal(bytes Bytes, v interface{}) error {
	switch v := v.(type) {
	case *Bytes:
		v.Copy(bytes)
		return nil
	default:
		return ErrLazyWrongValueType
	}
}
