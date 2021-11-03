package encoding

import (
	"encoding/xml"
)

// import()

type XML struct{}

func init() {
	register(new(XML))
}

func (*XML) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (*XML) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
