package encoding

import (
	"encoding/json"
	"errors"
)

var (
	ErrJSONWrongValueType = errors.New("codec JSON converts on wrong type value")
)

type JSON struct{}

func init() {
	register(new(JSON))
}

func (*JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (*JSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
