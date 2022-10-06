package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/s3d"
	"github.com/xackery/quail/wld"
)

var (
	isFound     bool
	nameRegex   = regexp.MustCompile(`IT([0-9]+).*`)
	name2Regex  = regexp.MustCompile(`it([0-9]+).*.mod`)
	name3Regex  = regexp.MustCompile(`it([0-9]+).*.mds`)
	filesWithIT = ""
	names       = make(map[string][]string)
)

func main() {
	start := time.Now()
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("ran in", time.Since(start).Seconds(), "seconds")
}

func run() error {

	path := "./"
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("entries: %w", err)
	}
	for _, entry := range entries {
		isFound = false
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".eqg" && ext != ".s3d" {
			continue
		}

		noExtName := entry.Name()
		noExtName = noExtName[0 : len(noExtName)-4]
		fmt.Println("parsing", entry.Name(), len(names), "entries so far")
		if ext == ".s3d" {
			err = parseS3D(path, noExtName, names)
			if err != nil {
				fmt.Println(noExtName, err)
			}
			continue
		}
		if ext == ".eqg" {
			err = parseEQG(path, noExtName, names)
			if err != nil {
				fmt.Println(noExtName, err)
			}
		}
	}

	w, err := os.Create("weapons.txt")
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer w.Close()

	w.WriteString(fmt.Sprintf("found at: %s\n", filesWithIT[0:len(filesWithIT)-2]))
	for name, entries := range names {
		w.WriteString(fmt.Sprintf("IT" + name + ": "))
		for _, zone := range entries {
			w.WriteString(zone + " ")
		}
		w.WriteString("\n")
	}
	return nil
}

func parseS3D(path string, name string, names map[string][]string) error {
	archive, err := s3d.NewFile(path + name + ".s3d")
	if err != nil {
		return fmt.Errorf("s3d newFile: %w", err)
	}
	e, err := wld.NewFile(name, archive, name+".wld")
	if err != nil {
		return fmt.Errorf("newfile wld: %w", err)
	}

	for _, name := range e.NameCache {
		if !strings.HasPrefix(name, "IT") {
			continue
		}
		if strings.Contains(name, "_") {
			continue
		}
		//fmt.Println(name)
		matches := nameRegex.FindAllStringSubmatch(name, -1)
		if len(matches) > 0 {
			add(name, matches[0][1])
			if !isFound {
				isFound = true
				filesWithIT += name + ".s3d, "
			}
		}
	}
	return nil
}

func parseEQG(path string, name string, names map[string][]string) error {
	archive, err := eqg.NewFile(path + name + ".eqg")
	if err != nil {
		return fmt.Errorf("eqg newFile: %w", err)
	}

	for _, entry := range archive.Files() {
		name := strings.ToLower(entry.Name())
		if !strings.HasPrefix(name, "it") {
			continue
		}
		//fmt.Println(name)
		matches := name2Regex.FindAllStringSubmatch(name, -1)
		if len(matches) > 0 {
			add(name, matches[0][1])
			if !isFound {
				isFound = true
				filesWithIT += name + ".eqg, "
				continue
			}
		}
		matches = name3Regex.FindAllStringSubmatch(name, -1)
		if len(matches) > 0 {
			add(name, matches[0][1])
			if !isFound {
				isFound = true
				filesWithIT += name + ".eqg, "
				continue
			}
		}
	}
	return nil
}

func add(key string, value string) {
	_, ok := names[key]
	if ok {
		names[key] = append(names[key], value)
	}
	names[key] = []string{value}
}
