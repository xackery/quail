package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed: ", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		fmt.Println("usage: audiodump <path>")
		os.Exit(1)
	}
	path := os.Args[1]
	fmt.Println("path:", path)

	log.SetLogLevel(1)

	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".pfs":
		default:
			return nil
		}

		pfs, err := pfs.NewFile(path)
		if err != nil {
			return fmt.Errorf("load: %w", err)
		}

		for _, fe := range pfs.Files() {

			if !strings.HasSuffix(fe.Name(), ".wav") {
				continue
			}

			fmt.Printf("%s|%s\n", fe.Name(), filepath.Base(path))
		}

		return nil
	})
}
