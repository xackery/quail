package main

import (
	"fmt"
	"time"

	"github.com/xackery/quail/cmd"
	"github.com/xackery/quail/log"
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

	if Version == "" {
		Version = fmt.Sprintf("dev-%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
		//fmt.Printf("quail %s\n", Version)
	}
	if ShowVersion == "1" {
		fmt.Printf("quail %s\n", Version)
	}
	log.SetLogLevel(1)
	//log.LogToFile()
	cmd.Execute()
}
