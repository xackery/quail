package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/pfs/s3d"
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Compress an eqg/s3d folder named _file.ext/ to a pfs archive",
	Long:  `Compress is used to take a provided _file.eqg or _file.s3d and compress it based on a folder structure`,
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
		err = compress(path, out)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)
	compressCmd.PersistentFlags().String("path", "", "path to compress")
	compressCmd.PersistentFlags().String("out", "", "name of compressed eqg archive output, defaults to path's basename")
	compressCmd.Example = `quail compress --path="./_clz.eqg/"
quail compress ./_soldungb.eqg/
quail compress _soldungb.eqg/ common.eqg
quail compress --path=_soldungb.eqg/ --out=foo.eqg`
}

func compress(path string, out string) error {
	if strings.HasSuffix(out, ".eqg") {
		return compressEQG(path, out)
	}
	if strings.HasSuffix(out, ".s3d") {
		return compressS3D(path, out)
	}
	if strings.HasSuffix(out, ".pfs") {
		return compressEQG(path, out)
	}

	out = out + ".eqg"
	return compressEQG(path, out)
}

func compressEQG(path string, out string) error {
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

		data, err := os.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name(), err)
		}
		err = archive.Add(file.Name(), data)
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
	err = archive.Encode(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", out, err)
	}

	fmt.Printf("%s\n%d file%s written to %s\n", addStdout, fileCount, helper.Pluralize(fileCount), out)
	return nil
}

func compressS3D(path string, out string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("path invalid, must be a directory (%s)", path)
	}

	archive := &s3d.S3D{}
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
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name(), err)
		}
		err = archive.Add(file.Name(), data)
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
	err = archive.Encode(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", out, err)
	}

	fmt.Printf("%s\n%d file%s written to %s\n", addStdout, fileCount, helper.Pluralize(fileCount), out)
	return nil
}
