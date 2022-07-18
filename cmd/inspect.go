package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/eqg"
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

		err = inspect(path)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.PersistentFlags().String("path", "", "path to inspect")

}

func inspect(path string) error {
	var err error
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".eqg":
		err = inspectEQG(path)
	default:
		return fmt.Errorf("unsupported extension: %s", ext)
	}
	if err != nil {
		return fmt.Errorf("inspectEQG: %w", err)
	}
	return nil
}

func inspectEQG(path string) error {
	archive, err := eqg.New(filepath.Base(path))
	if err != nil {
		return fmt.Errorf("new: %w", err)
	}
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()
	err = archive.Load(r)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	fmt.Printf("%s contains %d files:\n", filepath.Base(path), archive.Len())

	filesByName := archive.Files()
	sort.Sort(common.FilerByName(filesByName))
	for _, fe := range archive.Files() {
		base := float64(len(fe.Data()))
		out := ""
		num := float64(1024)
		if base < num*num*num*num {
			out = fmt.Sprintf("%0.0fG", base/num/num/num)
		}
		if base < num*num*num {
			out = fmt.Sprintf("%0.0fM", base/num/num)
		}
		if base < num*num {
			out = fmt.Sprintf("%0.0fK", base/num)
		}
		if base < num {
			out = fmt.Sprintf("%0.0fB", base)
		}
		fmt.Printf("%s\t%s\n", out, fe.Name())
	}

	return nil
}
