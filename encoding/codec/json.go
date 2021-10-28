package codec

import (
	stdjson "encoding/json"
	"errors"
	"reflect"
)

var (
	ErrJSONWrongValueType = errors.New("codec JSON converts on wrong type value")
)

type JSON struct{}

func init() {
	Register(new(JSON))
}

func (*JSON) String() string {
	return "json"
}

func (*JSON) Marshal(v interface{}) (Bytes, error) {
	value := reflect.ValueOf(v)
	if !(value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct || value.Kind() == reflect.Struct) {
		return nilBytes, ErrJSONWrongValueType
	}
	data, err := stdjson.Marshal(v)
	if err != nil {
		return nilBytes, err
	}
	return Bytes{data}, nil
}

func (*JSON) Unmarshal(bytes Bytes, v interface{}) error {
	value := reflect.ValueOf(v)
	if !(value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
		return ErrJSONWrongValueType
	}
	return stdjson.Unmarshal(bytes.Data, v)
}
