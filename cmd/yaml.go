package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"gopkg.in/yaml.v3"
)

// yamlCmd represents the yaml command
var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "Convert eq types to or from yaml",
	Long:  `Allows you to convert binary eq types to yaml, and visa versa`,
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("parse path: %w", err)
		}
		if srcPath == "" {
			if len(args) < 1 {
				return cmd.Usage()
			}
			srcPath = args[0]
		}
		dstPath, _ := cmd.Flags().GetString("file")
		if dstPath == "" {
			if len(args) >= 2 {
				dstPath = args[1]
			}
		}

		srcFile := ""
		dstFile := ""
		if strings.Contains(srcPath, ":") {
			srcFile = strings.Split(srcPath, ":")[1]
			srcPath = strings.Split(srcPath, ":")[0]
		}
		if strings.Contains(dstPath, ":") {
			dstFile = strings.Split(dstPath, ":")[1]
			dstPath = strings.Split(dstPath, ":")[0]
		}

		if srcPath == "" {
			return fmt.Errorf("src path is required")
		}

		if dstPath == "" {
			return fmt.Errorf("dst path is required")
		}

		fi, err := os.Stat(dstPath)
		if err == nil && fi.IsDir() {

			data, err := os.ReadFile(srcPath)
			if err != nil {
				return fmt.Errorf("read src file: %w", err)
			}

			type metaPeek struct {
				MetaFileName string `yaml:"file_name"`
			}

			meta := &metaPeek{}
			err = yaml.Unmarshal(data, meta)
			if err != nil {
				return fmt.Errorf("yaml unmarshal: %w", err)
			}
			dstPath = filepath.Join(dstPath, meta.MetaFileName)
		}

		isSrcArchive := false
		srcPathExt := strings.ToLower(filepath.Ext(srcPath))
		isDstArchive := false
		dstPathExt := strings.ToLower(filepath.Ext(dstPath))
		switch srcPathExt {
		case ".s3d", ".pfs", ".pak", ".eqg":
			isSrcArchive = true
		}
		switch dstPathExt {
		case ".s3d", ".pfs", ".pak", ".eqg":
			isDstArchive = true
		}

		if isSrcArchive && isDstArchive {
			return fmt.Errorf("source and destination cannot both be archives")
		}

		if isSrcArchive {
			return readYamlArchiveFile(srcPath, srcFile, dstPath)
		}
		if isDstArchive {
			return writeYamlArchiveFile(srcPath, dstPath, dstFile)
		}

		if dstPathExt == ".yaml" {
			return writeYamlFile(srcPath, dstPath)
		}
		return readYamlFile(srcPath, dstPath)
	},
}

func init() {
	rootCmd.AddCommand(yamlCmd)
	yamlCmd.PersistentFlags().String("path", "", "path to inspect")
	yamlCmd.PersistentFlags().String("file", "", "file to inspect inside pfs")
}

func writeYamlArchiveFile(srcYamlPath string, dstArchivePath string, dstArchiveFile string) error {
	if srcYamlPath == "" {
		return fmt.Errorf("src yaml path is required")
	}

	srcExt := filepath.Ext(srcYamlPath)

	dstExt := filepath.Ext(dstArchivePath)
	if dstExt != ".s3d" && dstExt != ".pfs" && dstExt != ".pak" && dstExt != ".eqg" {
		return fmt.Errorf("dst file must be archive")
	}

	archive, err := pfs.NewFile(dstArchivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}

	reader := raw.New(srcExt)
	if reader == nil {
		return fmt.Errorf("unsupported file format %s", srcExt)
	}

	buf := &bytes.Buffer{}
	err = yaml.NewDecoder(buf).Decode(reader)
	if err != nil {
		return fmt.Errorf("yaml decode: %w", err)
	}

	err = archive.Add(dstArchiveFile, buf.Bytes())
	if err != nil {
		return fmt.Errorf("add archive file: %w", err)
	}

	w, err := os.Create(dstArchivePath)
	if err != nil {
		return fmt.Errorf("create dst file: %w", err)
	}

	err = archive.Write(w)
	if err != nil {
		return fmt.Errorf("write archive: %w", err)
	}

	return nil
}

func writeYamlFile(srcYamlPath string, dstPath string) error {
	if srcYamlPath == "" {
		return fmt.Errorf("src yaml path is required")
	}

	srcExt := filepath.Ext(srcYamlPath)

	dstExt := filepath.Ext(dstPath)
	if dstExt != ".yaml" {
		return fmt.Errorf("dst file must be yaml")
	}

	r, err := os.Open(srcYamlPath)
	if err != nil {
		return fmt.Errorf("open src file: %w", err)
	}
	defer r.Close()

	reader := raw.New(srcExt)
	if reader == nil {
		return fmt.Errorf("unsupported file format %s", srcExt)
	}

	err = yaml.NewDecoder(r).Decode(reader)
	if err != nil {
		return fmt.Errorf("yaml decode: %w", err)
	}

	w, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create dst file: %w", err)
	}
	defer w.Close()

	err = reader.Write(w)
	if err != nil {
		return fmt.Errorf("write dst file: %w", err)
	}
	return nil

}

func readYamlArchiveFile(srcArchivePath string, srcArchiveFile string, dstYamlPath string) error {
	if srcArchivePath == "" {
		return fmt.Errorf("src archive path is required")
	}
	if srcArchiveFile == "" {
		return fmt.Errorf("src archive file is required")
	}

	srcExt := filepath.Ext(srcArchivePath)

	dstExt := filepath.Ext(dstYamlPath)
	if dstExt != ".yaml" {
		return fmt.Errorf("dst file must be yaml")
	}

	archive, err := pfs.NewFile(srcArchivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer archive.Close()

	data, err := archive.File(srcArchiveFile)
	if err != nil {
		return fmt.Errorf("read archive file: %w", err)
	}

	reader, err := raw.Read(srcExt, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("read archive file: %w", err)
	}

	w, err := os.Create(dstYamlPath)
	if err != nil {
		return fmt.Errorf("create dst file: %w", err)
	}
	defer w.Close()

	err = yaml.NewEncoder(w).Encode(reader)
	if err != nil {
		return fmt.Errorf("yaml encode: %w", err)
	}
	return nil
}

func readYamlFile(srcPath string, dstYamlPath string) error {
	if srcPath == "" {
		return fmt.Errorf("src path is required")
	}

	srcExt := filepath.Ext(srcPath)

	dstExt := filepath.Ext(dstYamlPath)
	if dstExt != ".yaml" {
		return fmt.Errorf("dst file must be yaml")
	}

	r, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open src file: %w", err)
	}
	defer r.Close()

	reader, err := raw.Read(srcExt, r)
	if err != nil {
		return fmt.Errorf("read src file: %w", err)
	}

	w, err := os.Create(dstYamlPath)
	if err != nil {
		return fmt.Errorf("create dst file: %w", err)
	}
	defer w.Close()

	err = yaml.NewEncoder(w).Encode(reader)
	if err != nil {
		return fmt.Errorf("yaml encode: %w", err)
	}

	return nil

}
