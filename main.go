package main

import (
	"fmt"
	"os"
	"time"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/s3d"
)

var (
	Version = "v0.0.0"
)

func main() {
	start := time.Now()
	err := run()
	if err != nil {
		log.Println("failed:", err.Error())
		os.Exit(1)
	}
	log.Printf("finished in %0.2f seconds\n", time.Since(start).Seconds())
}

func run() error {
	log.SetLogLevel(1)
	log.Println("quail", Version)
	path := "s3d/test/clz.wld"
	log.Println("working on", path)
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	e := &s3d.Wld{}
	err = e.Load(f)
	if err != nil {
		return fmt.Errorf("load: %v", err)
	}
	return nil
}
