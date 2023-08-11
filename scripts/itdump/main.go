package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/log"
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
		case ".eqg":
		case ".s3d":
		default:
			return nil
		}

		//fmt.Println(filepath.Base(path))
		q := quail.New()

		err = q.PFSImport(path)
		if err != nil {
			fmt.Println(filepath.Base(path), "pfs import:", err)
			return nil
			//return fmt.Errorf("pfs import: %w", err)
		}

		for _, mesh := range q.Meshes {
			if !strings.HasPrefix(strings.ToLower(mesh.Name), "it") {
				continue
			}
			fmt.Printf("%s|%s\n", mesh.Name, filepath.Base(path))
		}

		return nil
	})
}
