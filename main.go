package main

import (
	"fmt"

	"github.com/xackery/quail/cmd"
	"github.com/xackery/quail/common"
)

var (
	Version string
)

func main() {
	common.Version = Version
	if common.Version != "" {
		fmt.Printf("quail v%s\n", common.Version)
	}
	cmd.Execute()
}
