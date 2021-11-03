package encoding

import "errors"

var (
	ErrEmptyCodecCalled = errors.New("codec empty should not be called")
)

type Empty struct{}

func init() {
	registerCodec(new(Empty))
}

func (*Empty) Marshal(_ interface{}) (Bytes, error) {
	return nilBytes, ErrEmptyCodecCalled
}

func (s *Empty) Unmarshal(_ Bytes, _ interface{}) error {
	return ErrEmptyCodecCalled
}
