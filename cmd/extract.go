package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/pfs/s3d"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract an eqg or s3d archive to a _file.ext/ folder",
	Long:  `Extract an eqg or s3d archive`,
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
				return err
			}
			archiveCount := 0
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				ok, err := extract(fmt.Sprintf("%s/%s", path, file.Name()), fmt.Sprintf("%s/_%s", out, file.Name()), true)
				if err != nil {
					return err
				}
				if !ok {
					continue
				}
				archiveCount++
			}
			if archiveCount == 0 {
				return fmt.Errorf("no archives found")
			}
			fmt.Println("extracted", archiveCount, "archives")
			return nil
		}

		_, err = extract(path, out, false)
		if err != nil {
			return err
		}
		return nil
	},
}

func extract(in string, out string, isDir bool) (bool, error) {
	var err error
	switch strings.ToLower(filepath.Ext(in)) {
	case ".eqg":
		err = extractEQG(in, out, isDir)
		if err != nil {
			return true, fmt.Errorf("extractEQG: %w", err)
		}
	case ".s3d":
		err = extractS3D(in, out, isDir)
		if err != nil {
			return true, fmt.Errorf("extractS3D: %w", err)
		}
	case ".pfs":
		err = extractEQG(in, out, isDir)
		if err != nil {
			return true, fmt.Errorf("extractEQG: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}

func extractEQG(in string, out string, isDir bool) error {
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	a, err := eqg.New("out")
	if err != nil {
		return fmt.Errorf("eqg.New: %w", err)
	}
	err = a.Decode(f)
	if err != nil {
		return fmt.Errorf("decode %s: %w", in, err)
	}
	results, err := a.Extract(out)
	if err != nil {
		fmt.Printf("Failed to extract %s: %s\n", filepath.Base(in), err)
		os.Exit(1)
	}

	if isDir && len(results) > 64 {
		results = results[0:64] + " (summarized due to directory extract)"
	}
	fmt.Println(filepath.Base(in), results)
	return nil
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().String("path", "", "path to compressed eqg")
	extractCmd.PersistentFlags().String("out", "", "out path to extract to")
}

func extractS3D(path string, out string, isDir bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	e, err := s3d.New("out")
	if err != nil {
		return fmt.Errorf("eqg.New: %w", err)
	}
	err = e.Decode(f)
	if err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}

	results, err := e.Extract(out)
	if err != nil {
		fmt.Printf("Failed to extract %s: %s\n", filepath.Base(path), err)
		os.Exit(1)
	}

	if isDir && len(results) > 64 {
		results = results[0:64] + " (summarized due to directory extract)"
	}
	fmt.Println(filepath.Base(path), results)
	return nil
}
