package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a file",
	Long: `Inspect an EverQuest asset to discover contents within

Supported extensions: eqg, zon, ter, ani, mod
`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer f.Close()
		ext := strings.ToLower(filepath.Ext(path))

		//shortname := filepath.Base(path)
		//shortname = strings.TrimSuffix(shortname, filepath.Ext(shortname))
		ok, err := inspectEQG(f, ext)
		if err != nil {
			return fmt.Errorf("inspectEQG: %w", err)
		}
		if ok {
			return nil
		}

		return fmt.Errorf("failed to inspect: unknown extension %s on file %s", ext, filepath.Base(path))
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.PersistentFlags().String("path", "", "path to inspect")

}

func inspectEQG(f io.Reader, ext string) (bool, error) {
	if ext != ".eqg" {
		return false, nil
	}

	return true, nil
}
