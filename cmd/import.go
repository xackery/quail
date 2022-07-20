package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/s3d"
	"github.com/xackery/quail/ter"
	"github.com/xackery/quail/zon"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Export an eqg or s3d archive to embedded GLTF",
	Long:  `Export an eqg or s3d archive to embedded GLTF`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("parse path: %w", err)
		}
		if path == "" {
			if len(args) < 1 {
				return cmd.Usage()
			}
			path = args[0]
		}
		defer func() {
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}()
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("parse absolute path: %w", err)
		}
		if out == "" {
			if len(args) < 2 {
				out = filepath.Base(absPath)
			} else {
				out = args[1]
			}
		}

		out = strings.ToLower(out)

		if strings.Contains(out, ".") && !strings.HasSuffix(out, ".eqg") && !strings.HasSuffix(out, ".s3d") {
			return fmt.Errorf("only eqg and s3d extension out names are supported")
		}
		out = strings.TrimPrefix(out, "_")
		err = importPath(path, out)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().String("path", "", "path to import")
	importCmd.PersistentFlags().String("out", "", "name of imported eqg archive output, defaults to path's basename")
	importCmd.Example = `quail import --path="./_clz.eqg/"
quail import ./_soldungb.eqg/
quail import _soldungb.eqg/ foo.eqg
quail import --path=_soldungb.eqg/ --out=foo.eqg`
}

func importPath(path string, out string) error {
	if strings.HasSuffix(out, ".eqg") {
		return importEQG(path, out)
	}
	if strings.HasSuffix(out, ".s3d") {
		return importS3D(path, out)
	}

	out = out + ".eqg"
	return importEQG(path, out)
}

func importEQG(path string, out string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("path invalid, must be a directory (%s)", path)
	}

	archive := &eqg.EQG{}
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("readdir path: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found in %s to add to archive %s", path, out)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == ".DS_Store" {
			continue
		}
		fileName := strings.ToLower(file.Name())
		baseName := filepath.Base(fileName)
		if strings.Contains(baseName, ".") {
			baseName = baseName[0:strings.Index(baseName, ".")]
		}

		if strings.HasSuffix(fileName, ".zon.gltf") {
			e, err := zon.New(baseName, archive)
			if err != nil {
				return fmt.Errorf("zon new %s: %w", baseName, err)
			}

			gdoc, err := gltf.Open(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				return fmt.Errorf("zon open %s: %w", baseName, err)
			}
			err = e.GLTFImport(gdoc)
			if err != nil {
				return fmt.Errorf("zon import %s: %w", baseName, err)
			}

			err = e.ArchiveExport(archive)
			if err != nil {
				return fmt.Errorf("zon archive export %s: %w", baseName, err)
			}

			continue
		}

		if strings.HasSuffix(fileName, ".ter.gltf") {
			e, err := ter.New(baseName, archive)
			if err != nil {
				return fmt.Errorf("ter new %s: %w", baseName, err)
			}

			err = e.ArchiveExport(archive)
			if err != nil {
				return fmt.Errorf("ter archive export %s: %w", baseName, err)
			}

			continue
		}

		if strings.HasSuffix(fileName, ".mds.gltf") {
			e, err := mds.New(baseName, archive)
			if err != nil {
				return fmt.Errorf("mds new %s: %w", baseName, err)
			}

			err = e.ArchiveExport(archive)
			if err != nil {
				return fmt.Errorf("mds archive export %s: %w", baseName, err)
			}

			continue
		}

		if strings.HasSuffix(fileName, ".mod.gltf") {
			e, err := mds.New(baseName, archive)
			if err != nil {
				return fmt.Errorf("mod new %s: %w", baseName, err)
			}

			err = e.ArchiveExport(archive)
			if err != nil {
				return fmt.Errorf("mod archive export %s: %w", baseName, err)
			}

			continue
		}
	}

	w, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("create %s: %w", out, err)
	}
	defer w.Close()
	err = archive.Save(w)
	if err != nil {
		return fmt.Errorf("save %s: %w", out, err)
	}

	fmt.Printf("%d file%s written to %s\n", archive.Len(), helper.Pluralize(archive.Len()), out)
	return nil
}

func importS3D(path string, out string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("path invalid, must be a directory (%s)", path)
	}

	a := &s3d.S3D{}
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("readdir path: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found in %s to add to archive %s", path, out)
	}

	addStdout := ""
	fileCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == ".DS_Store" {
			continue
		}
		fileCount++
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name(), err)
		}
		err = a.Add(file.Name(), data)
		if err != nil {
			return fmt.Errorf("add %s: %w", file.Name(), err)
		}
		addStdout += file.Name() + ", "
	}
	if fileCount == 0 {
		return fmt.Errorf("no files found to add")
	}
	addStdout = addStdout[0:len(addStdout)-2] + "\n"

	w, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("create %s: %w", out, err)
	}
	defer w.Close()
	err = a.Save(w)
	if err != nil {
		return fmt.Errorf("save %s: %w", out, err)
	}

	fmt.Printf("%d file%s: %s\nwritten to %s\n", fileCount, helper.Pluralize(fileCount), addStdout, out)
	return nil
}
