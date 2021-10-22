package cluster

import "errors"

type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type bytesCodec struct{}

var BytesCodec bytesCodec

var (
	ErrMashallingInvalidType   = errors.New("invalid type mashalling with bytes codec")
	ErrUnmashallingInvalidType = errors.New("invalid type unmashalling with bytes codec")
)

func (bytesCodec) Marshal(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		copyData := make([]byte, len(data))
		copy(copyData, data)
		return copyData, nil
	}
	return nil, ErrMashallingInvalidType
}

func (bytesCodec) Unmarshal(data []byte, v interface{}) error {
	if _, ok := v.(*[]byte); ok {
		copyData := make([]byte, len(data))
		copy(copyData, data)
		v = &copyData
		return nil
	}
	return ErrUnmashallingInvalidType
}
