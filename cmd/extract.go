package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/eqg"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract an eqg archive to provided directory",
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
			if len(args) > 0 {
				path = args[0]
			} else {
				return cmd.Usage()
			}
		}
		path = strings.TrimSuffix(path, "/")
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}
		if out == "" {
			if len(args) > 1 {
				out = args[1]
			} else {
				out = fmt.Sprintf("./_%s", filepath.Base(path))
			}
		}

		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path check: %w", err)
		}
		if fi.IsDir() {
			if strings.Contains(out, "_") {
				out = fmt.Sprintf("./%s", filepath.Base(path))
			}
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				if strings.ToLower(filepath.Ext(file.Name())) != ".eqg" {
					continue
				}
				err = extract(fmt.Sprintf("%s/%s", path, file.Name()), fmt.Sprintf("%s/_%s", out, file.Name()), true)
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
			}
			os.Exit(0)
		}
		err = extract(path, out, false)
		if err != nil {
			fmt.Println("Error: extract:", err)
			os.Exit(1)
		}

		return nil
	},
}

func extract(path string, out string, isDir bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	e := &eqg.EQG{}
	err = e.Load(f)
	if err != nil {
		return fmt.Errorf("load %s: %w", path, err)
	}
	results, err := e.Extract(out)
	if err != nil {
		fmt.Printf("Error: extract %s: %s\n", filepath.Base(path), err)
		os.Exit(1)
	}

	if isDir && len(results) > 64 {
		results = results[0:64] + " (summarized due to directory extract)"
	}
	fmt.Println(filepath.Base(path), results)
	return nil
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().String("path", "", "path to compressed eqg")
	extractCmd.PersistentFlags().String("out", "", "out folder to extract to")
}
