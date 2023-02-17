package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/model/metadata/zon"
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug a file",
	Long: `Debug an EverQuest asset to discover contents within

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
				out = fmt.Sprintf("debug_%s.png", filepath.Base(path))
			} else {
				out = args[1]
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
			return fmt.Errorf("debug requires a target file, directory provided")
		}

		dump.New(path)
		defer dump.Close()
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open debug path: %s", err)
		}
		defer f.Close()
		ext := strings.ToLower(filepath.Ext(path))

		//shortname := filepath.Base(path)
		//shortname = strings.TrimSuffix(shortname, filepath.Ext(shortname))
		type decoder interface {
			Decode(io.ReadSeeker) error
		}
		type decodeTypes struct {
			instance  decoder
			extension string
		}
		decodes := []*decodeTypes{
			//{instance: &ani.ANI{}, extension: ".ani"},
			//{instance: &eqg.EQG{}, extension: ".eqg"},
			//{instance: &mod.MOD{}, extension: ".mod"},
			//{instance: &ter.TER{}, extension: ".ter"},
			{instance: &zon.ZON{}, extension: ".zon"},
		}

		for _, v := range decodes {
			if ext != v.extension {
				continue
			}

			err = v.instance.Decode(f)
			if err != nil {
				return fmt.Errorf("failed to decode %s: %w", v.extension, err)
			}
			return nil
		}
		return fmt.Errorf("failed to debug: unknown extension %s on file %s", ext, filepath.Base(path))
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.PersistentFlags().String("path", "", "path to debug")
	debugCmd.PersistentFlags().String("out", "", "out file of debug")

}
