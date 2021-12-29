package encoding

import (
	"encoding/json"

	"github.com/aceaura/libra/core/magic"
)

type Hash struct{}

func init() {
	register(NewHash())
}

func NewHash() *Hash {
	return new(Hash)
}

func (h Hash) String() string {
	return magic.TypeName(h)
}

func (h Hash) Style() EncodingStyleType {
	return EncodingStyleStruct
}

func (Hash) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (Hash) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (h Hash) Reverse() Encoding {
	return h
}
