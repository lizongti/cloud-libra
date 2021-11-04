package encoding

import (
	"encoding/xml"

	"github.com/aceaura/libra/magic"
)

type XML struct{}

func init() {
	register(new(XML))
}

func NewXML() *XML {
	return new(XML)
}

func (xml XML) String() string {
	return magic.TypeName(xml)
}

func (XML) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (XML) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
