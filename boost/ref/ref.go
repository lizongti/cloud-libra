package ref

import "reflect"

func TypeName(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		return reflect.TypeOf(i).Elem().Name()
	} else if v.Kind() == reflect.Struct {
		return reflect.TypeOf(i).Name()
	}
	return ""
}

func CallName(t interface{}, method string) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		reflect.TypeOf(i).Elem().MethodByName(method).Call([]reflect.Value{})
	} else if v.Kind() == reflect.Struct {
		reflect.TypeOf(i).MethodByName(method).Call([]reflect.Value{})
	}
}
