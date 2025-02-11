package main

import (
	"fmt"
	"time"

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

	info, _ := debug.ReadBuildInfo()
	Version = info.Main.Version
	if Version == "" {
		Version = "dev-" + time.Now().Format("20060102")
	}

	fmt.Printf("Quail version %s, Go version: %s\n", Version, info.GoVersion)
	helper.Version = Version
	cmd.Execute()
}
