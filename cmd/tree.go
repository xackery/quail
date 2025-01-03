/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/quail"
)

func init() {
	rootCmd.AddCommand(treeCmd)
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
	if len(args) < 1 {
		fmt.Println("Usage: quail tree <src>")
		os.Exit(1)
	}
	start := time.Now()
	defer func() {
		fmt.Printf("Finished in %0.2f seconds\n", time.Since(start).Seconds())
	}()
	srcPath := args[0]
	fi, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("tree: srcPath is a directory. Set to a file for this extension")
	}
	q := quail.New()

	err = q.TreeRead(srcPath)
	if err != nil {
		return fmt.Errorf("tree read: %w", err)
	}

	return nil
}
