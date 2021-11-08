package encoding

import (
	"encoding/json"
	"errors"

	"github.com/aceaura/libra/magic"
)

var (
	ErrJSONWrongValueType = errors.New("encoding JSON converts on wrong type value")
)

type JSON struct{}

func init() {
	register(NewJSON())
}

func NewJSON() *JSON {
	return new(JSON)
}

func (json JSON) String() string {
	return magic.TypeName(json)
}

func (JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j JSON) Reverse() Encoding {
	return j
}
