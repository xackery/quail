package main

import (
	"fmt"
	"time"

	"github.com/xackery/quail/cmd"
	"github.com/xackery/quail/log"
)

var (
	// Version is the current version
	Version string
)

func main() {
	if Version == "" {
		Version = fmt.Sprintf("dev-%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
		//fmt.Printf("quail %s\n", Version)
	}
	log.SetLogLevel(1)
	//log.LogToFile()
	cmd.Execute()
}
