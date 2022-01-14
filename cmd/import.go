/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		cachePath := fmt.Sprintf("%s/cache/", path)
		fi, err = os.Stat(cachePath)
		if err != nil {
			fmt.Printf("cache folder not found for import: %s\n", err)
			os.Exit(1)
		}
		if !fi.IsDir() {
			fmt.Println("cache file exists, should be folder")
			os.Exit(1)
		}

		e := &eqg.EQG{}
		files, err := os.ReadDir(cachePath)
		if err != nil {
			return fmt.Errorf("readdir cachepath: %w", err)
		}
		if len(files) == 0 {
			fmt.Printf("no files found in %s to import\n", cachePath)
			os.Exit(1)
		}

		addStdout := ""
		fileCount := 0
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileCount++

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
		if fileCount == 0 {
			fmt.Println("no files found to add")
			os.Exit(1)
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

		fmt.Printf("%s created with %d file%s: %s", out, fileCount, helper.Pluralize(fileCount), addStdout)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().String("path", "", "path to import")
	importCmd.PersistentFlags().String("out", "", "name of compressed eqg archive output, defaults to path's basename")
	importCmd.Example = `quail import --path="./_clz.eqg/"`
}
