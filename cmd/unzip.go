package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/pfs"
)

func init() {
	rootCmd.AddCommand(unzipCmd)
	unzipCmd.PersistentFlags().String("path", "", "path to compressed pfs")
	unzipCmd.PersistentFlags().String("out", "", "out path to unzip to")

}

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "Unzip a pfs archive (eqg, s3d, pak, pfs)",
	Long:  `Unzip a pfs archive to either provided out folder or to _file.ext/ folder`,
	Example: `quail unzip clz.eqg # unzips to ./_clz.eqg/
quail unzip foo.eqg foofolder/
quail unzip --path=foo.eqg --out=foofolder/"`,
	Run: runUnzip,
}

func runUnzip(cmd *cobra.Command, args []string) {
	err := runUnzipE(cmd, args)
	if err != nil {
		log.Printf("Failed: %s", err.Error())
		os.Exit(1)
	}
}

func runUnzipE(cmd *cobra.Command, args []string) error {
	var srcArchivePath string
	var err error

	if cmd != nil {
		srcArchivePath, err = cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("parse path: %w", err)
		}
	}
	if srcArchivePath == "" {
		if len(args) > 0 {
			srcArchivePath = args[0]
		} else {
			return cmd.Usage()
		}
	}

	srcArchivePath = strings.TrimSuffix(srcArchivePath, "/")
	srcFile := ""

	if strings.Contains(srcArchivePath, ":") {
		srcFile = strings.Split(srcArchivePath, ":")[1]
		srcArchivePath = strings.Split(srcArchivePath, ":")[0]
	}

	var dstPath string
	if cmd != nil {
		dstPath, err = cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}
	}
	if dstPath == "" && srcFile == "" {
		dstPath = fmt.Sprintf("./_%s", filepath.Base(srcArchivePath))
		if len(args) > 1 {
			dstPath = args[1]
		}
	}

	if srcFile != "" {
		archive, err := pfs.NewFile(srcArchivePath)
		if err != nil {
			return fmt.Errorf("pfs.NewFile: %w", err)
		}

		if dstPath == "" {
			dstPath = "."
		}

		di, err := os.Stat(dstPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("dst path check: %w", err)
		}
		if di != nil && di.IsDir() {
			dstPath = filepath.Join(dstPath, srcFile)
		}

		log.Printf("Unziping %s:%s to %s", srcArchivePath, srcFile, dstPath)
		data, err := archive.File(srcFile)
		if err != nil {
			return fmt.Errorf("archive.File: %w", err)
		}
		err = os.WriteFile(dstPath, data, 0644)
		if err != nil {
			return fmt.Errorf("write %s: %w", dstPath, err)
		}
		return nil
	}

	fi, err := os.Stat(srcArchivePath)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if fi.IsDir() {
		dstPath = filepath.Join(dstPath, "_"+filepath.Base(srcArchivePath))
	}

	log.Printf("Unziping %s to %s", srcArchivePath, dstPath)
	err = os.MkdirAll(dstPath, 0755)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	archive, err := pfs.NewFile(srcArchivePath)
	if err != nil {
		return fmt.Errorf("pfs.NewFile: %w", err)
	}

	fileCount := 0
	for _, fe := range archive.Files() {
		fePath := filepath.Join(dstPath, fe.Name())
		err = os.WriteFile(fePath, fe.Data(), 0644)
		if err != nil {
			return fmt.Errorf("write %s: %w", fePath, err)
		}
		fileCount++
	}
	log.Printf("Unzipped %d files", fileCount)

	return nil
}
