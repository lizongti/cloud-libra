package encoding

import "errors"

var (
	ErrEmptyEncodingCalled = errors.New("emtpy encoding should not be called")
)

type EmtpyEncoding struct{}

var emptyEncoding = new(EmtpyEncoding)

func Emtpy() TypeEncoding {
	return emptyEncoding
}

func (*EmtpyEncoding) String() string {
	return "empty"
}

func (*EmtpyEncoding) Marshal(_ interface{}) ([]byte, error) {
	return nil, ErrEmptyEncodingCalled
}

func (s *EmtpyEncoding) Unmarshal(_ []byte, _ interface{}) error {
	return ErrEmptyEncodingCalled
}
