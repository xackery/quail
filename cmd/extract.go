package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
			return cmd.Usage()
		}
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}
		if out == "" {
			out = fmt.Sprintf("./_%s", filepath.Base(path))
		}

		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path check: %w", err)
		}
		if fi.IsDir() {
			return fmt.Errorf("inspect requires a target file, directory provided")
		}

		f, err := os.Open(path)
		if err != nil {
			fmt.Println("Error: open:", err)
			os.Exit(1)
		}
		e := &eqg.EQG{}
		err = e.Load(f)
		if err != nil {
			fmt.Printf("Error: load %s: %s\n", filepath.Base(path), err)
			os.Exit(1)
		}
		err = e.Extract(out)
		if err != nil {
			fmt.Printf("Error: extract %s: %s\n", filepath.Base(path), err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().String("path", "", "path to compressed eqg")
	extractCmd.PersistentFlags().String("out", "", "out folder to extract to")
}
