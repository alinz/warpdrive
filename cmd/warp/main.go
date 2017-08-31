package main

import (
	"fmt"
	"os"

	"github.com/pressly/warpdrive/cmd/warp/section"
)

func main() {
	if err := section.Root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
