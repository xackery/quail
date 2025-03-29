package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		fmt.Println("usage: zondump <path>")
		fmt.Println("path points to an eq dir, generates a shaderdump.txt with leaked data")
		os.Exit(1)
	}
	path := os.Args[1]
	fmt.Println("path:", path)

	zones := make(map[string]int)

	start := time.Now()

	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		baseName := filepath.Base(path)
		zonName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		switch ext {
		case ".eqg":
		case ".zon":
			r, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open %s: %w", zonName, err)
			}
			zon, err := raw.Read(".zon", r)
			if err != nil {
				return fmt.Errorf("zon read %s: %w", zonName, err)
			}
			switch z := zon.(type) {
			case *raw.Zon:
				fmt.Println(zonName, z.Version)
				zones[zonName] = int(z.Version)
			default:
				return fmt.Errorf("unknown type %T", z)
			}
			return nil
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

		fmt.Println(baseName)

		for _, file := range a.Files() {
			fileExt := filepath.Ext(file.Name())
			switch fileExt {
			case ".zon":

			default:
				continue
			}

			r := bytes.NewReader(file.Data())

			rawRead, err := raw.Read(fileExt, r)
			if err != nil {
				fmt.Printf("raw read %s: %v\n", filepath.Base(file.Name()), err)
				continue
			}
			switch dat := rawRead.(type) {
			case *raw.Zon:
				fmt.Println("> ", zonName, dat.Version)

				zones[zonName] = int(dat.Version)
			default:
				return fmt.Errorf("unknown type %T", dat)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}

	w, err := os.Create("zondump.txt")
	if err != nil {
		return fmt.Errorf("create zondump.txt: %w", err)
	}
	defer w.Close()

	// Extract keys from the map
	zoneNames := make([]string, 0, len(zones))
	for name := range zones {
		zoneNames = append(zoneNames, name)
	}

	// Sort the keys alphabetically
	sort.Strings(zoneNames)

	// Iterate through sorted keys and write to file
	for _, name := range zoneNames {
		version := zones[name]
		_, err := w.WriteString(fmt.Sprintf("%s|%d\n", name, version))
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}

	fmt.Printf("Finished %0.2f seconds\n", time.Since(start).Seconds())
	return nil
}
