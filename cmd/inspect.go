package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/ani"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
	"github.com/xackery/quail/zon"
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
		ext := strings.ToLower(filepath.Ext(path))

		shortname := filepath.Base(path)
		shortname = strings.TrimSuffix(shortname, filepath.Ext(shortname))
		type loader interface {
			Load(io.ReadSeeker, string) error
		}
		type loadTypes struct {
			instance  loader
			extension string
		}
		loads := []*loadTypes{
			{instance: &ani.ANI{}, extension: ".ani"},
			{instance: &eqg.EQG{}, extension: ".eqg"},
			{instance: &mod.MOD{}, extension: ".mod"},
			{instance: &ter.TER{}, extension: ".ter"},
			{instance: &zon.ZON{}, extension: ".zon"},
		}

		for _, v := range loads {
			if ext != v.extension {
				continue
			}

			err = v.instance.Load(f, shortname)
			if err != nil {
				fmt.Printf("Error: load %s: %s\n", v.extension, err)
				os.Exit(1)
			}
			return nil
		}

		fmt.Printf("Error: inspect: unknown extension %s on file %s\n", ext, filepath.Base(path))
		os.Exit(1)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.PersistentFlags().String("path", "", "path to inspect")
	inspectCmd.PersistentFlags().String("out", "", "out file of inspect")

}
