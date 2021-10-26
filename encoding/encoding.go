package encoding

import (
	"errors"
)

type Encoding interface {
	String() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

var (
	ErrEmptyEncodingCalled = errors.New("emtpy encoding should not be called")
	ErrWrongValueType      = errors.New("encoding convert on wrong type value")
)

func Marshal(c Encoding, v interface{}) ([]byte, error) {
	return c.Marshal(v)
}

func Unmarshal(c Encoding, data []byte, v interface{}) error {
	return c.Unmarshal(data, v)
}
