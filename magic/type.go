package magic

import (
	"context"
	"reflect"
)

var (
	TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	TypeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	TypeOfBytes   = reflect.TypeOf(([]byte)(nil))
	TypeNil       = reflect.Type(nil)
)

func TypeName(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		return reflect.TypeOf(i).Elem().Name()
	} else if v.Kind() == reflect.Struct {
		return reflect.TypeOf(i).Name()
	}
	return ""
}
