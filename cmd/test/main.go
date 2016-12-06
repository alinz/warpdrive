package main

import "fmt"
import "strings"

type Header struct {
	filename [1000]byte
}

func main() {
	header := &Header{}

	src := []byte("./ali/test")

	for i := 0; i < len(src); i++ {
		header.filename[i] = src[i]
	}

	path := string(header.filename[:])
	path = strings.TrimSpace(path)

	fmt.Println(len(path))
}
