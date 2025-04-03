package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/quail"
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
		fmt.Println("usage: modelfind <name> <path>")
		fmt.Println("path points to an eq dir, generates a texturefind.txt with leaked data")
		os.Exit(1)
	}
	name := os.Args[1]
	path := os.Args[2]
	fmt.Println("path:", path)

	start := time.Now()

	w, err := os.Create("modelfind.txt")
	if err != nil {
		return fmt.Errorf("create modelfind.txt: %w", err)
	}
	defer w.Close()

	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		baseName := filepath.Base(path)

		ext := filepath.Ext(path)
		switch ext {
		case ".eqg":
			fmt.Println(baseName)
			a, err := pfs.NewFile(path)
			for _, file := range a.Files() {
				fileExt := filepath.Ext(file.Name())
				switch fileExt {
				case ".mod":
				case ".mds":

				default:
					continue
				}
				if !strings.Contains(file.Name(), name) {
					continue
				}
				fmt.Printf("%s: %s", baseName, file.Name())
				fmt.Fprintf(w, "%s: %s\n", baseName, file.Name())
			}
			if err != nil {
				fmt.Printf("pfs open %s: %v\n", filepath.Base(path), err)
				return nil
			}

			return nil

		case ".s3d":
		//case ".pak":
		default:
			return nil
		}

		fmt.Println(baseName)
		q := quail.New()

		err = q.PfsRead(path)
		if err != nil {
			fmt.Println(filepath.Base(path), "pfs import:", err)
			return nil
			//return fmt.Errorf("pfs import: %w", err)
		}

		if q.Wld == nil {
			return nil
		}

		for _, actor := range q.Wld.ActorDefs {
			if strings.Contains(actor.Tag, name) {
				fmt.Printf("%s: %s", baseName, actor.Tag)
				fmt.Fprintf(w, "%s: %s\n", baseName, actor.Tag)
			}
		}

		for _, sprite := range q.Wld.DMSpriteDef2s {
			if strings.Contains(sprite.Tag, name) {
				fmt.Printf("%s: %s", baseName, sprite.Tag)
				fmt.Fprintf(w, "%s: %s\n", baseName, sprite.Tag)
			}
		}
		for _, anim := range q.Wld.TrackDefs {
			if strings.Contains(anim.Tag, name) {
				fmt.Printf("%s: %s", baseName, anim.Tag)
				fmt.Fprintf(w, "%s: %s\n", baseName, anim.Tag)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walkdir: %w", err)
	}

	fmt.Printf("Finished %0.2f seconds\n", time.Since(start).Seconds())
	return nil
}
