package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"gopkg.in/yaml.v3"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect an EverQuest asset",
	Long: `Inspect an EverQuest asset to discover contents within
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
		file, err := cmd.Flags().GetString("file")
		if strings.Contains(path, ":") {
			file = strings.Split(path, ":")[1]
			path = strings.Split(path, ":")[0]
		}

		if file == "" {
			if len(args) >= 2 {
				file = args[1]
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
			return fmt.Errorf("inspect requires a target file, directory provided")
		}

		inspected, err := inspect(path, file)
		if err != nil {
			return fmt.Errorf("inspect: %w", err)
		}

		buf := &bytes.Buffer{}
		enc := yaml.NewEncoder(buf)
		enc.SetIndent(2)
		err = enc.Encode(inspected)
		if err != nil {
			return fmt.Errorf("yaml encode: %w", err)
		}

		log.Infoln(buf.String())

		//reflectTraversal(inspected, 0, -1)
		return nil

	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.PersistentFlags().String("path", "", "path to inspect")
	inspectCmd.PersistentFlags().String("file", "", "file to inspect inside pfs")
}

func inspect(path string, file string) (interface{}, error) {
	/* if len(file) < 2 {
		//log.Infof("Inspecting %s", filepath.Base(path))
	} else {
		//log.Infof("Inspecting %s %s", filepath.Base(path), filepath.Base(file))
	} */

	isValidExt := false
	exts := []string{".eqg", ".s3d", ".pfs", ".pak"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return inspectFile(nil, path, file)
	}

	pfs, err := pfs.NewFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s load: %w", ext, err)
	}
	if len(file) < 2 {
		pfs.Close()
		return *pfs, nil
	}
	return inspectFile(pfs, path, file)
}

func inspectFile(pfs *pfs.Pfs, path string, file string) (interface{}, error) {
	if pfs == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return inspectContent(filepath.Base(path), bytes.NewReader(data))
	}

	for _, fe := range pfs.Files() {
		if !strings.EqualFold(fe.Name(), file) {
			continue
		}
		return inspectContent(file, bytes.NewReader(fe.Data()))
	}
	return nil, fmt.Errorf("%s not found in %s", file, filepath.Base(path))
}

func inspectContent(file string, r *bytes.Reader) (interface{}, error) {
	ext := strings.ToLower(filepath.Ext(file))
	reader, err := raw.Read(ext, r)
	if err != nil {
		return nil, fmt.Errorf("%s read: %w", ext, err)
	}
	reader.SetFileName(filepath.Base(file))
	return reader, nil
}
