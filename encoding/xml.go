package encoding

import (
	"encoding/xml"
	"reflect"
)

// import()

type XML struct{}

func init() {
	registerCodec(new(XML))
}

func (*XML) Marshal(v interface{}) (Bytes, error) {
	value := reflect.ValueOf(v)
	if !(value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct || value.Kind() == reflect.Struct) {
		return nilBytes, ErrJSONWrongValueType
	}
	data, err := xml.Marshal(v)
	if err != nil {
		return nilBytes, err
	}
	return MakeBytes(data), nil
}

func (*XML) Unmarshal(bytes Bytes, v interface{}) error {
	value := reflect.ValueOf(v)
	if !(value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
		return ErrJSONWrongValueType
	}
	return xml.Unmarshal(bytes.Data, v)
}
