package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/mod"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Review an EverQuest file",
	Long: `Review an EverQuest file and inspect it's contents.

Supported extensions:
- eqg: everquest archive
- zon: zone definition
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("parse path: %w", err)
		}
		if path == "" {
			return cmd.Usage()
		}

		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path check: %w", err)
		}
		if fi.IsDir() {
			return fmt.Errorf("inspect requires a target file, directory provided")
		}

		fmt.Println("inspect: generated ./inspect.png")
		d, err := dump.New(filepath.Base(path))
		if err != nil {
			return fmt.Errorf("dump.New: %w", err)
		}
		defer d.Save("inspect.png")
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("Error: open:", err)
			os.Exit(1)
		}
		defer f.Close()
		ext := filepath.Ext(path)
		switch strings.ToLower(ext) {
		case ".eqg":
			e := &eqg.EQG{}
			err = e.Load(f)
			if err != nil {
				fmt.Printf("Error: load %s: %s\n", filepath.Base(path), err)
				return nil
			}
		case ".mod":
			e := &mod.MOD{}
			err = e.Load(f)
			if err != nil {
				fmt.Printf("Error: load %s: %s\n", filepath.Base(path), err)
				return nil
			}
		default:
			fmt.Printf("Error: inspect: unknown extension %s on path %s", ext, filepath.Base(path))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.PersistentFlags().String("path", "", "path to inspect")

}
