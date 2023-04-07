package main

import (
	"fmt"

	"github.com/xackery/quail/cmd"
)

var (
	Version string
)

func main() {
	fmt.Printf("quail v%s\n", Version)
	cmd.Execute()
}
