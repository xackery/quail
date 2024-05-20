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
		fmt.Println("usage: leak <path>")
		fmt.Println("path points to an eq dir, generates a leaks.txt with leaked data")
		os.Exit(1)
	}
	path := os.Args[1]
	fmt.Println("path:", path)

	log.SetLogLevel(1)

	w, err := os.Create("leaks.txt")
	if err != nil {
		return fmt.Errorf("create leaks.txt: %w", err)
	}
	defer w.Close()

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
		case ".pak":
		default:
			return nil
		}

		a, err := pfs.NewFile(path)
		if err != nil {
			fmt.Printf("pfs open %s: %v\n", filepath.Base(path), err)
			return nil
		}

		allowedExt := []string{".exe", ".sph", ".spk", ".sps", ".mdf", ".spr", ".spk"}
		for _, file := range a.Files() {

			fileExt := filepath.Ext(file.Name())
			if !contains(allowedExt, fileExt) {
				continue
			}

			w.WriteString(fmt.Sprintf("%s %s\n", file.Name(), filepath.Base(path)))
		}
		return nil
	})
}

func contains(arr []string, str string) bool {
	str = strings.ToLower(str)
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
