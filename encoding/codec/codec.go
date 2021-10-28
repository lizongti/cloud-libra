package codec

type Bytes struct {
	Data []byte
}

func (b Bytes) Dulplicate() Bytes {
	data := make([]byte, len(b.Data))
	copy(data, b.Data)
	return Bytes{data}
}

func (b Bytes) Copy(in Bytes) {
	b.Data = make([]byte, len(in.Data))
	copy(b.Data, in.Data)
}

var nilBytes = Bytes{}

type Codec interface {
	String() string
	Marshal(interface{}) (Bytes, error)
	Unmarshal(Bytes, interface{}) error
}

func Marshal(e Codec, v interface{}) ([]byte, error) {
	bytes, err := e.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.Data, nil
}

func Unmarshal(e Codec, data []byte, v interface{}) error {
	return e.Unmarshal(Bytes{data}, v)
}

func MarshalBytes(e Codec, v interface{}) (Bytes, error) {
	return e.Marshal(v)
}

func UnmarshalBytes(e Codec, bytes Bytes, v interface{}) error {
	return e.Unmarshal(bytes, v)
}

type CodecSet map[string]Codec

var codecSet = make(map[string]Codec)

func Register(codecs ...Codec) {
	for _, codec := range codecs {
		codecSet[codec.String()] = codec
	}
}
