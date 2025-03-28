package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func init() {
	rootCmd.AddCommand(zipCmd)
	zipCmd.PersistentFlags().String("path", "", "path to zip")
	zipCmd.PersistentFlags().String("out", "", "name of zipped eqg archive output, defaults to path's basename")
	zipCmd.Example = `quail zip --path="./_clz.eqg/"
quail zip ./_soldungb.eqg/
quail zip _soldungb.eqg/ helper.eqg
quail zip --path=_soldungb.eqg/ --out=foo.eqg`
}

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Zip a folder to a pfs archive (eqg, s3d, pfs or pak)",
	Long: `Zip is used to take a provided folder and zip it.
	There is a shorthand system where if you only provide a folder with no destination, it will use the folder's name to determine the output file.`,
	Example: `quail zip _clz.eqg
quail zip somefolder somepfs.s3d
quail zip --path=somefolder --out=somepfs.s3d
`,
	Run: runZip,
}

func runZip(cmd *cobra.Command, args []string) {
	err := runZipE(cmd, args)
	if err != nil {
		log.Printf("Failed: %s", err.Error())
		os.Exit(1)
	}
}

func runZipE(cmd *cobra.Command, args []string) error {
	var path string
	var err error
	if cmd != nil {
		path, err = cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("parse path: %w", err)
		}
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
	var out string
	if cmd != nil {
		out, err = cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}
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

	isValid := false
	for _, ext := range []string{".eqg", ".s3d", ".pfs", ".pak"} {
		if strings.HasSuffix(out, ext) {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("out must have a valid extension (.eqg, .s3d, .pfs, .pak)")
	}
	out = strings.TrimPrefix(out, "_")
	err = zip(path, out)
	if err != nil {
		return err
	}
	return nil
}

func zip(path string, out string) error {
	if strings.HasSuffix(out, ".eqg") {
		return zipPfs(path, out)
	}
	if strings.HasSuffix(out, ".s3d") {
		return zipPfs(path, out)
	}
	if strings.HasSuffix(out, ".pfs") {
		return zipPfs(path, out)
	}
	if strings.HasSuffix(out, ".pak") {
		return zipPfs(path, out)
	}

	out = out + ".eqg"
	return zipPfs(path, out)
}

func zipPfs(path string, out string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("path invalid, must be a directory (%s)", path)
	}

	archive := &pfs.Pfs{}
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
	err = archive.Write(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", out, err)
	}

	fmt.Printf("%s\n%d file%s written to %s\n", addStdout, fileCount, helper.Pluralize(fileCount), out)
	return nil
}
