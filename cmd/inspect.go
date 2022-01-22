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
	"github.com/xackery/quail/ter"
	"github.com/xackery/quail/zon"
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
			if len(args) < 1 {
				return cmd.Usage()
			}
			path = args[0]
		}
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}
		if out == "" {
			if len(args) < 2 {
				out = "inspect.png"
			} else {
				out = args[1]
			}
		}

		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path check: %w", err)
		}
		if fi.IsDir() {
			return fmt.Errorf("inspect requires a target file, directory provided")
		}

		fmt.Println("inspect: generated", out)
		d, err := dump.New(filepath.Base(path))
		if err != nil {
			return fmt.Errorf("dump.New: %w", err)
		}
		defer d.Save(out)
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
		case ".ter":
			e := &ter.TER{}
			err = e.Load(f)
			if err != nil {
				fmt.Printf("Error: load %s: %s\n", filepath.Base(path), err)
				return nil
			}
		case ".zon":
			e := &zon.ZON{}
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
	inspectCmd.PersistentFlags().String("out", "", "out file of inspect")

}
