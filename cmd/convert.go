package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/quail"
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

	if srcExt == ".quail" {
		err = q.DirRead(srcPath)
		if err != nil {
			return fmt.Errorf("quail read dir: %w", err)
		}
	} else {
		err = q.PfsRead(srcPath)
		if err != nil {
			return fmt.Errorf("pfs read: %w", err)
		}
	}

	dstExt := filepath.Ext(dstPath)
	if dstExt == ".quail" {
		err = q.DirWrite(dstPath)
		if err != nil {
			return fmt.Errorf("dir write: %w", err)
		}
		return nil
	}

	err = q.PfsWrite(1, 1, dstPath)
	if err != nil {
		return fmt.Errorf("pfs write: %w", err)
	}

	return nil
}
