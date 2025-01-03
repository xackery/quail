package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"os"

	"github.com/xackery/quail/quail"
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
		fmt.Println("usage: itdump <path>")
		os.Exit(1)
	}
	path := os.Args[1]
	fmt.Println("path:", path)

	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".eqg":
		case ".s3d":
		default:
			return nil
		}

		//fmt.Println(filepath.Base(path))
		q := quail.New()

		err = q.PfsRead(path)
		if err != nil {
			fmt.Println(filepath.Base(path), "pfs import:", err)
			return nil
			//return fmt.Errorf("pfs import: %w", err)
		}

		for _, model := range q.Models {
			if !strings.HasPrefix(strings.ToLower(model.Header.Name), "it") {
				continue
			}
			fmt.Printf("%s|%s\n", model.Header.Name, filepath.Base(path))
		}

		return nil
	})
}
