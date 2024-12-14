package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/os"
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
		fmt.Println("usage: itdump <path>")
		os.Exit(1)
	}
	path := os.Args[1]
	fmt.Println("path:", path)

	var prefixAnims = make(map[string]bool)
	var suffixAnims = make(map[string]bool)

	aw, err := os.Create("all.txt")
	if err != nil {
		return fmt.Errorf("create all.txt: %w", err)
	}
	defer aw.Close()

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
		case ".s3d":
			return nil
		default:
			return nil
		}
		pfs, err := pfs.NewFile(path)
		if err != nil {
			return fmt.Errorf("eqg new: %w", err)
		}

		pfsName := filepath.Base(path)
		for _, file := range pfs.Files() {
			name := file.Name()
			fext := filepath.Ext(name)
			switch fext {
			case ".ani":
				animName := name[0 : len(name)-4]
				_, err = aw.WriteString(fmt.Sprintf("%s|%s\n", pfsName, animName))
				if err != nil {
					return fmt.Errorf("write all.txt: %w", err)
				}
				if strings.Contains(animName, "_") {
					records := strings.Split(animName, "_")
					prefixAnims[records[0]] = true
					suffixAnims[records[len(records)-1]] = true
				}
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}

	pw, err := os.Create("prefix.txt")
	if err != nil {
		return fmt.Errorf("create prefix.txt: %w", err)
	}
	defer pw.Close()
	for anim := range prefixAnims {
		pw.WriteString(anim + "\n")
	}

	sw, err := os.Create("suffix.txt")
	if err != nil {
		return fmt.Errorf("create suffix.txt: %w", err)
	}
	defer sw.Close()
	for anim := range suffixAnims {
		sw.WriteString(anim + "\n")
	}
	return nil
}
