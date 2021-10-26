package encoding

import "encoding/json"

type JSONEncoding struct{}

var jsonEncoding = new(JSONEncoding)

func JSON() Encoding {
	return jsonEncoding
}

func (*JSONEncoding) String() string {
	return "json"
}

func (*JSONEncoding) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s *JSONEncoding) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
