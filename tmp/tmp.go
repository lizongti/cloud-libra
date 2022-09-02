package main

import (
	"fmt"

	. "github.com/go-python/cpy3"
)

func main() {
	Py_Initialize()
	gostr := "foo"
	bytes := PyBytes_FromString(gostr)
	str := PyBytes_AsString(bytes)
	fmt.Println("hello [", str, "]")
}
