package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/helper"
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Create an eqg archive by compressing a directory",
	Long:  `Compress is used to compress an eqg archive based on provided arguments`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("parse path: %w", err)
		}
		if path == "" {
			return cmd.Usage()
		}
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("parse absolute path: %w", err)
		}
		if out == "" {
			out = filepath.Base(absPath)
		}

		out = strings.ToLower(out)

		if strings.Contains(out, ".") && !strings.HasSuffix(out, ".eqg") {
			return fmt.Errorf("only .eqg extension out names are supported")
		}

		if !strings.HasSuffix(out, ".eqg") {
			out = out + ".eqg"
		}
		out = strings.TrimPrefix(out, "_")

		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path check: %w", err)
		}
		if !fi.IsDir() {
			return fmt.Errorf("path invalid, must be a directory (%s)", path)
		}

		e := &eqg.EQG{}
		files, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("readdir path: %w", err)
		}
		if len(files) == 0 {
			return fmt.Errorf("no files found in %s to add to archive %s", path, out)
		}

		addStdout := ""
		for _, file := range files {

			data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				return fmt.Errorf("read %s: %w", file.Name(), err)
			}
			err = e.Add(file.Name(), data)
			if err != nil {
				return fmt.Errorf("add %s: %w", file.Name(), err)
			}
			addStdout += file.Name() + ", "
		}
		addStdout = addStdout[0:len(addStdout)-2] + "\n"

		w, err := os.Create(out)
		if err != nil {
			return fmt.Errorf("create %s: %w", out, err)
		}
		defer w.Close()
		err = e.Save(w)
		if err != nil {
			return fmt.Errorf("save %s: %w", out, err)
		}

		fmt.Printf("%s created with %d file%s: %s", out, len(files), helper.Pluralize(len(files)), addStdout)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)
	compressCmd.PersistentFlags().String("path", "", "path to compress")
	compressCmd.PersistentFlags().String("out", "", "name of compressed eqg archive output, defaults to path's basename")
	compressCmd.Example = `quail compress --path="./_clz.eqg/"`
}
