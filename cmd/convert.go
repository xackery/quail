package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/quail"
	"github.com/xackery/quail/raw"
)

func init() {
	rootCmd.AddCommand(convertCmd)
}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Take one file, convert to another",
	Long: `Supports eqg, s3d, and quail (wcemu) files
Usage: quail convert <src> <dst>
Example: quail convert foo.s3d foo.quail - Takes foo.s3d and creates a folder called foo.quail
Example: quail convert foo.quail foo.s3d - Takes foo.quail folder and creates a foo.s3d file`,
	RunE: runConvert,
}

func runConvert(cmd *cobra.Command, args []string) error {
	err := runConvertE(cmd, args)
	if err != nil {
		fmt.Printf("Failed: %s\n", err.Error())
		os.Exit(1)
	}
	return nil
}

func runConvertE(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: quail convert <src> <dst>")
		os.Exit(1)
	}

	start := time.Now()
	defer func() {
		fmt.Printf("Finished in %0.2f seconds\n", time.Since(start).Seconds())
	}()
	srcPath := args[0]
	dstPath := args[1]
	fi, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}
	srcExt := filepath.Ext(srcPath)
	if fi.IsDir() && srcExt != ".quail" {
		return fmt.Errorf("convert: srcPath is %s but also a directory. Set to a file for this extension", srcExt)
	}

	q := quail.New()

	switch srcExt {
	case ".quail":
		err = q.DirRead(srcPath)
		if err != nil {
			return fmt.Errorf("quail read dir: %w", err)
		}

	case ".json":
		err = q.JsonRead(srcPath)
		if err != nil {
			return fmt.Errorf("json read: %w", err)
		}
	default:
		baseName := filepath.Base(srcPath)
		err = q.PfsRead(srcPath)
		if err != nil {
			return fmt.Errorf("convert %s: %w", baseName, err)
		}
		if srcExt == ".eqg" {
			srcPathNoExt := srcPath[:len(srcPath)-len(srcExt)] // remove the .eqg
			err = quailLoadSideFile(q, srcPathNoExt+".zon")
			if err != nil {
				return fmt.Errorf("load side file .zon: %w", err)
			}
		}
	}

	dstExt := filepath.Ext(dstPath)
	switch dstExt {
	case ".quail":
		err = q.DirWrite(dstPath)
		if err != nil {
			return fmt.Errorf("dir write: %w", err)
		}

		err = os.MkdirAll(dstPath+"/.vscode/", 0755)
		if err != nil {
			return fmt.Errorf("mkdir vscode: %w", err)
		}

		ext := filepath.Ext(srcPath)
		w, err := os.Create(dstPath + "/.vscode/settings.json")
		if err != nil {
			return fmt.Errorf("create vscode settings: %w", err)
		}
		defer w.Close()
		w.WriteString("{\n")
		w.WriteString("    // quail path can be set via your environment PATH, a vscode workspace path, or uncomment below\n")
		w.WriteString("    // \"wce-vscode.quailPath\": \"/path/to/quail\"\n")
		w.WriteString("    \"wce-vscode.convertOnSave\": true,\n")
		w.WriteString(fmt.Sprintf("    \"wce-vscode.convertTargetPath\": \"temp%s\"\n", ext))
		w.WriteString("}\n")
		return nil
	case ".json":
		err = q.JsonWrite(dstPath)
		if err != nil {
			return fmt.Errorf("json write: %w", err)
		}
	default:
		err = q.PfsWrite(1, 1, dstPath)
		if err != nil {
			return fmt.Errorf("pfs write: %w", err)
		}

	}

	return nil
}

func quailLoadSideFile(q *quail.Quail, path string) error {
	ext := filepath.Ext(path)
	r, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // side files are optional
		}
		return fmt.Errorf("open side file %s: %w", ext, err)
	}
	defer r.Close()
	rawSideFile, err := raw.Read(ext, r) // read the raw data
	if err != nil {
		return fmt.Errorf("read %s side file: %w", ext, err)
	}
	if q.Wld == nil {
		return fmt.Errorf("no quail wld found to load side file %s", ext)
	}
	err = q.Wld.ReadRaw(rawSideFile) // read the raw data into Wld
	if err != nil {
		return fmt.Errorf("raw read side file %s: %w", ext, err)
	}
	fmt.Printf("Loaded side file %s\n", filepath.Base(path))
	return nil
}
