package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	fmt.Println("path:", path)

	outDir := filepath.Join(path, "textures")
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return fmt.Errorf("mkdir textures: %w", err)
	}

	err = os.MkdirAll(filepath.Join(outDir, "zones"), 0755)
	if err != nil {
		return fmt.Errorf("mkdir zones: %w", err)
	}

	err = os.MkdirAll(filepath.Join(outDir, "models"), 0755)
	if err != nil {
		return fmt.Errorf("mkdir models: %w", err)
	}

	textureCount := 0
	totalArchives := 0

	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		//case ".eqg":
		case ".s3d":
		default:
			return nil
		}

		totalArchives++
		if totalArchives%20 == 0 {
			fmt.Println("Archives processed:", totalArchives)
		}
		fmt.Println(filepath.Base(path))

		q := quail.New()
		defer q.Close()

		err = q.PfsRead(path)
		if err != nil {
			fmt.Println(filepath.Base(path), "pfs import:", err)
			return nil
			//return fmt.Errorf("pfs import: %w", err)
		}

		isModel := strings.Contains(path, "_chr")

		for name, texture := range q.Textures {

			outPath := ""
			if isModel {
				outPath = filepath.Join(outDir, "models", name)
			} else {
				outPath = filepath.Join(outDir, "zones", name)
			}
			_, err := os.Stat(outPath)
			if os.IsExist(err) {
				continue
			}

			err = os.WriteFile(outPath, texture, 0644)
			if err != nil {
				fmt.Printf("write texture %s: %s\n", name, err)
				continue
			}
			textureCount++
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}

	fmt.Println("Total archives processed:", totalArchives)
	fmt.Println("Total textures dumped:", textureCount)
	return nil
}
