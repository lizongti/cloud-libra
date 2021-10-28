package encoding

type Encoding struct {
	encode []string
	decode []string
}

func (e Encoding) Reverse() Encoding {
	return Encoding{
		encode: e.decode,
		decode: e.encode,
	}
}

func (e Encoding) Marshal(v interface{}) ([]byte, error) {
	return nil, nil
}

func (e Encoding) Unmarshal(data []byte, v interface{}) error {
	return nil
}

type EncodingSet map[string]Encoding

type EncodingType string

var encodingSet = make(map[string]Encoding)

func init() {}

type TypeEncoding interface {
	String() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type BinaryEncoding interface {
	String() string
	Marshal([]byte) ([]byte, error)
	Unmarshal([]byte) ([]byte, error)
}

func TypeMarshal(e TypeEncoding, v interface{}) ([]byte, error) {
	return e.Marshal(v)
}

func TypeUnmarshal(e TypeEncoding, data []byte, v interface{}) error {
	return e.Unmarshal(data, v)
}

func BinaryMarshal(e BinaryEncoding, data []byte) ([]byte, error) {
	return e.Marshal(data)
}

func BinaryUnmarshal(e BinaryEncoding, data []byte) ([]byte, error) {
	return e.Unmarshal(data)
}

type Bytes struct {
	data []byte
}

func (b Bytes) Data() []byte {
	return b.data
}
