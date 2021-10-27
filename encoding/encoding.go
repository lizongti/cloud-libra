package encoding

import (
	"errors"
)

var (
	ErrEmptyEncodingCalled = errors.New("emtpy encoding should not be called")
	ErrWrongValueType      = errors.New("encoding convert on wrong type value")
)

type TypeEncoding interface {
	String() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type BinaryEncoding interface {
	String() string
	MarshalBinary([]byte) ([]byte, error)
	UnmarshalBinary([]byte) ([]byte, error)
}

func TypeMarshal(e TypeEncoding, v interface{}) ([]byte, error) {
	return e.Marshal(v)
}

func TypeUnmarshal(e TypeEncoding, data []byte, v interface{}) error {
	return e.Unmarshal(data, v)
}

func BinaryMarshal(e BinaryEncoding, data []byte) ([]byte, error) {

}

func BinaryUnmarshal(e BinaryEncoding, data []byte) ([]byte, error) {

}
