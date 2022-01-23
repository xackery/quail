package main

import (
	"fmt"

	"github.com/xackery/quail/cmd"
)

var (
	Version string
)

func main() {
	if Version != "" {
		fmt.Printf("quail v%s\n", Version)
	}
	cmd.Execute()
}
