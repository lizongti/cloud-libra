package encoding

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrJSONWrongValueType = errors.New("codec JSON converts on wrong type value")
)

type JSON struct{}

func init() {
	registerCodec(new(JSON))
}

func (*JSON) Marshal(v interface{}) (Bytes, error) {
	value := reflect.ValueOf(v)
	if !(value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct || value.Kind() == reflect.Struct) {
		return nilBytes, ErrJSONWrongValueType
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nilBytes, err
	}
	return MakeBytes(data), nil
}

func (*JSON) Unmarshal(bytes Bytes, v interface{}) error {
	value := reflect.ValueOf(v)
	if !(value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
		return ErrJSONWrongValueType
	}
	return json.Unmarshal(bytes.Data, v)
}
