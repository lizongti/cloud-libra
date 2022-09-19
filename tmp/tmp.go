package main

import (
	"fmt"
	"os"
	// . "github.com/go-python/cpy3".
)

func main() {
	infos, err := os.ReadDir("tmp/tmp.go")
	fmt.Println(err)
	fmt.Println(os.IsNotExist(err))
	for _, info := range infos {
		fmt.Println(info.Name())
	}
}
