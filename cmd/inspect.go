package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.PersistentFlags().String("path", "", "path to inspect")
	inspectCmd.PersistentFlags().String("path2", "", "path to compare")
}

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect an EverQuest asset",
	Long:  `Inspect an EverQuest asset to discover contents within`,
	Run:   runInspect,
}

func runInspect(cmd *cobra.Command, args []string) {
	err := runInspectE(cmd, args)
	if err != nil {
		fmt.Printf("Failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func runInspectE(cmd *cobra.Command, args []string) error {
	var err error
	defer func() {
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}()

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

	path2, err := cmd.Flags().GetString("path2")
	if err != nil {
		return fmt.Errorf("parse path2: %w", err)
	}
	if path2 == "" && len(args) > 1 {
		path2 = args[1]
	}

	val, err := inspect(path)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	if path2 == "" {
		fmt.Printf("%#v\n", val)
		return nil
	}
	val2, err := inspect(path2)
	if err != nil {
		return fmt.Errorf("inspect2: %w", err)
	}
	return helper.Compare(val, val2)

}

func inspect(path string) (string, error) {
	var err error
	var file string
	var index int

	if strings.Contains(path, ":") {
		split := strings.Split(path, ":")
		if len(split) < 2 {
			return "", fmt.Errorf("invalid path")
		}
		path = split[0]
		file = split[1]
		if len(split) == 3 {
			index, err = strconv.Atoi(split[2])
			if err != nil {
				return "", fmt.Errorf("index: %w", err)
			}
		}
	}

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
		if len(file) > 0 {
			index, err = strconv.Atoi(file)
			if err != nil {
				return "", fmt.Errorf("index: %w", err)
			}
			file = ""
		}
		return inspectFile(nil, path, file, index)
	}

	pfs, err := pfs.NewFile(path)
	if err != nil {
		return "", fmt.Errorf("%s load: %w", ext, err)
	}
	if len(file) < 2 {
		pfs.Close()
		buf := new(bytes.Buffer)
		err = inspectDump(pfs, index, buf)
		if err != nil {
			return "", fmt.Errorf("dump: %w", err)
		}

		return buf.String(), nil
	}
	return inspectFile(pfs, path, file, index)
}

func inspectFile(pfs *pfs.Pfs, path string, file string, index int) (string, error) {
	if pfs == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return inspectContent(filepath.Base(path), bytes.NewReader(data), index)
	}

	for _, fe := range pfs.Files() {
		if !strings.EqualFold(fe.Name(), file) {
			continue
		}
		return inspectContent(file, bytes.NewReader(fe.Data()), index)
	}
	return "", fmt.Errorf("%s not found in %s", file, filepath.Base(path))
}

func inspectContent(file string, r *bytes.Reader, index int) (string, error) {
	ext := strings.ToLower(filepath.Ext(file))
	reader, err := raw.Read(ext, r)
	if err != nil {
		return "", fmt.Errorf("%s read: %w", ext, err)
	}

	reader.SetFileName(filepath.Base(file))
	buf := new(bytes.Buffer)
	err = inspectDump(reader, index, buf)
	if err != nil {
		return "", fmt.Errorf("dump: %w", err)
	}
	return buf.String(), nil
}

func inspectDump(inspect interface{}, index int, w io.Writer) error {
	switch val := inspect.(type) {
	case *raw.Wld:
		if index <= 0 {
			fmt.Fprintf(w, "%s\n", val.String())
			return nil
		}
		for i := 0; i < len(val.Fragments); i++ {
			if i+1 != index {
				continue
			}

			_, err := fmt.Fprintf(w, "%+v\n", val.Fragments[i])
			if err != nil {
				return fmt.Errorf("dump: %w", err)
			}
			return nil
		}
	default:
		if index > 0 {
			return fmt.Errorf("index not supported for %T", val)
		}
		_, err := fmt.Fprintf(w, "%#v\n", inspect)
		if err != nil {
			return fmt.Errorf("dump: %w", err)
		}
		return nil
	}
	return nil
}
