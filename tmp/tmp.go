package main

import (
	"fmt"
	"reflect"
)

type M struct {
}

func (M) Func1() error {
	return nil
}

func (*M) Func2() error {
	return nil
}

func (M) func3() error {
	return nil
}

func (*M) func4() error {
	return nil
}

func main() {
	t := reflect.TypeOf(M{})
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Println(method.Name)
	}

	t = reflect.TypeOf(&M{})
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Println(method.Name)
		method.Func.Call()
	}
}
