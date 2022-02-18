package ref_test

import (
	"fmt"
	"reflect"
	"testing"
)

type T struct{}

func (t *T) Foo() {
	fmt.Println("foo")
}

func TestCallName(t *testing.T) {
	var o T
	reflect.ValueOf(&o).MethodByName("Foo").Call([]reflect.Value{})
}
