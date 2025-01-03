package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(yamlCmd)
	yamlCmd.PersistentFlags().String("path", "", "path to inspect")
	yamlCmd.PersistentFlags().String("file", "", "file to inspect inside pfs")
}

// yamlCmd represents the yaml command
var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "Convert eq types to or from yaml",
	Long:  `Allows you to convert binary eq types to yaml, and visa versa`,
	Run:   runYaml,
}

func runYaml(cmd *cobra.Command, args []string) {
	err := runYamlE(cmd, args)
	if err != nil {
		log.Printf("Failed: %s", err.Error())
		os.Exit(1)
	}
}

func runYamlE(cmd *cobra.Command, args []string) error {
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
		err = archiveToYamlConvert(srcPath, srcFile, dstPath)
		if err != nil {
			return fmt.Errorf("eqg->yaml read: %w", err)
		}
		return nil
	}
	if isDstArchive {
		err = yamlToArchiveConvert(srcPath, dstPath, dstFile)
		if err != nil {
			return fmt.Errorf("yaml->eqg write: %w", err)
		}
		return nil
	}

	if dstPathExt == ".yaml" {
		err = fileToYamlConvert(srcPath, dstPath)
		if err != nil {
			return fmt.Errorf("eqFile->yaml read: %w", err)
		}
		return nil
	}
	err = yamlToFileConvert(srcPath, dstPath)
	if err != nil {
		return fmt.Errorf("yaml->eqFile write: %w", err)
	}
	return nil
}

func yamlToArchiveConvert(srcYamlPath string, dstArchivePath string, dstArchiveFile string) error {
	if srcYamlPath == "" {
		return fmt.Errorf("src yaml path is required")
	}

	//srcExt := filepath.Ext(srcYamlPath)

	dstArchiveExt := filepath.Ext(dstArchivePath)
	if dstArchiveExt != ".s3d" && dstArchiveExt != ".pfs" && dstArchiveExt != ".pak" && dstArchiveExt != ".eqg" {
		return fmt.Errorf("dst file must be archive")
	}

	dstExt := filepath.Ext(dstArchiveFile)

	log.Printf("Converting %s to %s:%s", srcYamlPath, dstArchivePath, dstArchiveFile)

	archive, err := pfs.NewFile(dstArchivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}

	reader := raw.New(dstExt)
	if reader == nil {
		return fmt.Errorf("unsupported file format %s", dstExt)
	}

	r, err := os.Open(srcYamlPath)
	if err != nil {
		return fmt.Errorf("open src file: %w", err)
	}
	defer r.Close()
	err = yaml.NewDecoder(r).Decode(reader)
	if err != nil {
		return fmt.Errorf("yaml decode: %w", err)
	}

	buf := &bytes.Buffer{}
	err = reader.Write(buf)
	if err != nil {
		return fmt.Errorf("write raw file %s: %w", dstArchiveFile, err)
	}

	err = archive.Set(dstArchiveFile, buf.Bytes())
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

func yamlToFileConvert(srcYamlPath string, dstPath string) error {
	if srcYamlPath == "" {
		return fmt.Errorf("src yaml path is required")
	}

	srcExt := filepath.Ext(srcYamlPath)
	if srcExt != ".yaml" {
		return fmt.Errorf("src file must be yaml")
	}
	dstExt := filepath.Ext(dstPath)

	log.Printf("Converting %s to %s", srcYamlPath, dstPath)

	r, err := os.Open(srcYamlPath)
	if err != nil {
		return fmt.Errorf("open src file: %w", err)
	}
	defer r.Close()

	reader := raw.New(dstExt)
	if reader == nil {
		return fmt.Errorf("unsupported file format %s", dstExt)
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

func archiveToYamlConvert(srcArchivePath string, srcArchiveFile string, dstYamlPath string) error {
	if srcArchivePath == "" {
		return fmt.Errorf("src archive path is required")
	}
	if srcArchiveFile == "" {
		return fmt.Errorf("src archive file is required")
	}

	srcExt := filepath.Ext(srcArchiveFile)

	dstExt := filepath.Ext(dstYamlPath)
	if dstExt != ".yaml" {
		return fmt.Errorf("dst file must be yaml")
	}
	log.Printf("Converting %s:%s to %s", srcArchivePath, srcArchiveFile, dstYamlPath)

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
		return fmt.Errorf("read raw file %s: %w", srcArchiveFile, err)
	}
	reader.SetFileName(srcArchiveFile)

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

func fileToYamlConvert(srcPath string, dstYamlPath string) error {
	if srcPath == "" {
		return fmt.Errorf("src path is required")
	}

	srcExt := filepath.Ext(srcPath)

	dstExt := filepath.Ext(dstYamlPath)
	if dstExt != ".yaml" {
		return fmt.Errorf("dst file must be yaml")
	}

	log.Printf("Converting %s to %s", srcPath, dstYamlPath)

	r, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open src file: %w", err)
	}
	defer r.Close()

	reader, err := raw.Read(srcExt, r)
	if err != nil {
		return fmt.Errorf("read src file: %w", err)
	}
	reader.SetFileName(filepath.Base(srcPath))

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
