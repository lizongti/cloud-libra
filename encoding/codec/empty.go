package codec

import "errors"

var (
	ErrEmptyCodecCalled = errors.New("codec empty should not be called")
)

type Empty struct{}

func init() {
	Register(new(Empty))
}

func (*Empty) String() string {
	return "empty"
}

func (*Empty) Marshal(_ interface{}) (Bytes, error) {
	return nilBytes, ErrEmptyCodecCalled
}

func (s *Empty) Unmarshal(_ Bytes, _ interface{}) error {
	return ErrEmptyCodecCalled
}
