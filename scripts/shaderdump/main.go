package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
		fmt.Println("usage: shaderdump <path>")
		fmt.Println("path points to an eq dir, generates a shaderdump.txt with leaked data")
		os.Exit(1)
	}
	path := os.Args[1]
	fmt.Println("path:", path)

	shaders := make(map[string]map[string]bool)

	start := time.Now()

	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
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
		fmt.Println("reading", baseName)

		for _, file := range a.Files() {
			fileExt := filepath.Ext(file.Name())
			switch fileExt {
			case ".mod":
			case ".mds":
			case ".ter":
			default:
				continue
			}

			fmt.Println("> ", file.Name())
			r := bytes.NewReader(file.Data())

			rawRead, err := raw.Read(fileExt, r)
			if err != nil {
				fmt.Printf("raw read %s: %v\n", filepath.Base(file.Name()), err)
				continue
			}
			switch dat := rawRead.(type) {
			case *raw.Mds:
				for _, mat := range dat.Materials {
					if shaders[mat.ShaderName] == nil {
						shaders[mat.ShaderName] = make(map[string]bool)
					}
					for _, prop := range mat.Properties {
						if shaders[mat.ShaderName][prop.Name] == false {
							shaders[mat.ShaderName][prop.Name] = true
						}
					}
				}
			case *raw.Mod:
				for _, mat := range dat.Materials {
					if shaders[mat.ShaderName] == nil {
						shaders[mat.ShaderName] = make(map[string]bool)
					}
					for _, prop := range mat.Properties {
						if shaders[mat.ShaderName][prop.Name] == false {
							shaders[mat.ShaderName][prop.Name] = true
						}
					}
				}
			case *raw.Ter:
				for _, mat := range dat.Materials {
					if shaders[mat.ShaderName] == nil {
						shaders[mat.ShaderName] = make(map[string]bool)
					}
					for _, prop := range mat.Properties {
						if shaders[mat.ShaderName][prop.Name] == false {
							shaders[mat.ShaderName][prop.Name] = true
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}

	w, err := os.Create("shaderdump.txt")
	if err != nil {
		return fmt.Errorf("create shaderdump.txt: %w", err)
	}
	defer w.Close()

	for shader, props := range shaders {
		_, err := w.WriteString(fmt.Sprintf("if shader == '%s' and (", shader))
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
		isFirst := true
		for prop := range props {
			if isFirst {
				_, err := w.WriteString(fmt.Sprintf("property == '%s' ", prop))
				if err != nil {
					return fmt.Errorf("write: %w", err)
				}
			} else {
				_, err := w.WriteString(fmt.Sprintf("or property == '%s' ", prop))
				if err != nil {
					return fmt.Errorf("write: %w", err)
				}
			}
		}
		_, err = w.WriteString("):\nreturn True\n")
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}

	fmt.Printf("Finished %0.2f seconds\n", time.Since(start).Seconds())
	return nil
}
