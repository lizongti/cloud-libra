package encoding


type Encoding struct {
	encode []Codec
	decode []Codec
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

