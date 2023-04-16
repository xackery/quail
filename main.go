package main

import (
	"fmt"

	"github.com/xackery/quail/cmd"
	"github.com/xackery/quail/log"
)

var (
	Version string
)

func main() {
	fmt.Printf("quail v%s\n", Version)
	log.SetLogLevel(1)
	log.LogToFile()
	cmd.Execute()
}
