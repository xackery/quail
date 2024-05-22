package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/xackery/quail/log"
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
	if len(os.Args) < 2 {
		fmt.Println("usage: leak <path>")
		fmt.Println("path points to an eq dir, generates a leaks.txt with leaked data")
		os.Exit(1)
	}
	path := os.Args[1]
	fmt.Println("path:", path)

	log.SetLogLevel(1)

	w, err := os.Create("rawdump.txt")
	if err != nil {
		return fmt.Errorf("create rawdump.txt: %w", err)
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

		for _, file := range a.Files() {
			fileExt := filepath.Ext(file.Name())
			if fileExt != ".wld" {
				continue
			}
			outPath := "_" + filepath.Base(path)

			err = os.MkdirAll(outPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("mkdir: %w", err)
			}
			outPath = filepath.Join(outPath, file.Name())
			err = os.WriteFile(outPath+".bin", file.Data(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("write file: %w", err)
			}

			cmd := exec.Command("wine", "wldcom-patch.exe", "-d", outPath+".bin", outPath)
			buf := bytes.Buffer{}
			cmd.Stdout = &buf

			err = cmd.Run()
			if err != nil {
				os.Remove(outPath + ".bin")
				w.WriteString(fmt.Sprintf("%s %s: %s %s\n", filepath.Base(path), file.Name(), err, buf.String()))
				continue
			}
			err = os.Remove(outPath + ".bin")
			if err != nil {
				return fmt.Errorf("remove: %w", err)
			}

		}
		return nil
	})
}
