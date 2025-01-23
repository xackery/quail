/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/quail"
)

func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.PersistentFlags().String("path", "", "path to pfs archive or file")
	treeCmd.PersistentFlags().String("path2", "", "path to compare pfs archive or file")
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
	var err error
	defer func() {
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}()

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

	path2, err := cmd.Flags().GetString("path2")
	if err != nil {
		return fmt.Errorf("parse path2: %w", err)
	}
	if path2 == "" {
		if len(args) >= 2 {
			path2 = args[1]
		}
	}

	q := quail.New()
	if path2 != "" {
		file1 := filepath.Base(path)
		file2 := filepath.Base(path2)
		fmt.Printf("Comparing %s to %s\n", file1, file2)
		err = q.TreeCompare(path, path2)
		if err != nil {
			return fmt.Errorf("tree compare: %w", err)
		}
		return nil
	}

	err = q.TreeRead(os.Stdout, path)
	if err != nil {
		return fmt.Errorf("tree read: %w", err)
	}

	return nil
}
