/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/quail"
)

func init() {
	rootCmd.AddCommand(convertCmd)
}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Take one file, convert to another",
	Long:  `Supports eqg, s3d, and quail files`,
	RunE:  runConvert,
}

func runConvert(cmd *cobra.Command, args []string) error {
	err := runConvertE(cmd, args)
	if err != nil {
		log.Printf("Failed: %s", err.Error())
		os.Exit(1)
	}
	return nil
}

func runConvertE(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: quail convert <src> <dst>")
		os.Exit(1)
	}
	log.SetLogLevel(0)
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

	err = q.PfsRead(srcPath)
	if err != nil {
		return fmt.Errorf("pfs read: %w", err)
	}

	err = q.PfsWrite(1, 1, dstPath)
	if err != nil {
		return fmt.Errorf("pfs write: %w", err)
	}

	return nil
}
