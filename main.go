package main

import (
	"fmt"

	"runtime/debug"

	"github.com/xackery/quail/cmd"
	"github.com/xackery/quail/helper"
	/* "net/http"         // part of pprof heap
	_ "net/http/pprof" // part of pprof heap */)

var (
	// Version is the current version
	Version string
	// ShowVersion is a flag to show version
	ShowVersion string
)

func main() {

	/* // pprof heap
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}() */

	info, ok := debug.ReadBuildInfo()
	if ok {
		Version = info.Main.Version
	}

	fmt.Printf("quail %s\n", Version)
	helper.Version = Version
	cmd.Execute()
}
