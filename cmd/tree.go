/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/quail"
)

func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.PersistentFlags().String("path", "", "path to pfs archive or file")
	treeCmd.PersistentFlags().String("file", "", "file to read inside pfs archive")

}

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Tree Node System",
	Long: `Creates a tree node visualization for files
Example: quail tree foo.s3d
`,
	RunE: runTree,
}

func runTree(cmd *cobra.Command, args []string) error {
	err := runTreeE(cmd, args)
	if err != nil {
		fmt.Printf("Failed: %s\n", err.Error())
		os.Exit(1)
	}
	return nil
}

func runTreeE(cmd *cobra.Command, args []string) error {
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
	file, err := cmd.Flags().GetString("file")
	if strings.Contains(path, ":") {
		file = strings.Split(path, ":")[1]
		path = strings.Split(path, ":")[0]
	}
	if file == "" {
		if len(args) >= 2 {
			file = args[1]
		}
	}
	defer func() {
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}()

	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("inspect requires a target file, directory provided")
	}
	q := quail.New()

	err = q.TreeRead(path, file)
	if err != nil {
		return fmt.Errorf("tree read: %w", err)
	}

	return nil
}
