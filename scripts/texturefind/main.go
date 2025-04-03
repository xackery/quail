package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xackery/quail/pfs"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 3 {
		fmt.Println("usage: texturefind <name> <path>")
		fmt.Println("path points to an eq dir, generates a texturefind.txt with leaked data")
		os.Exit(1)
	}
	name := os.Args[1]
	path := os.Args[2]
	fmt.Println("path:", path)

	start := time.Now()

	w, err := os.Create("texturefind.txt")
	if err != nil {
		return fmt.Errorf("create texturefind.txt: %w", err)
	}
	defer w.Close()

	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".eqg":
		//case ".s3d":
		//case ".pak":
		default:
			return nil
		}

		a, err := pfs.NewFile(path)
		if err != nil {
			fmt.Printf("pfs open %s: %v\n", filepath.Base(path), err)
			return nil
		}

		baseName := filepath.Base(path)

		for _, file := range a.Files() {
			fileExt := filepath.Ext(file.Name())
			switch fileExt {
			case ".dds":
			default:
				continue
			}
			if file.Name() != name {
				continue
			}

			fmt.Printf("%s:%s\n", baseName, file.Name())
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}

	fmt.Printf("Finished %0.2f seconds\n", time.Since(start).Seconds())
	return nil
}
