package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Reading.")
		c, err := in.ReadByte()
		if err == io.EOF {
			break
		}
		fmt.Print(string(c))
	}
}
