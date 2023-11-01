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

		if isSrcArchive && dstPathExt != ".yaml" {
			return fmt.Errorf("destination or source must be yaml")
		}

		if isDstArchive && srcPathExt != ".yaml" {
			return fmt.Errorf("destination or source must be yaml")
		}

		if dstPathExt == ".yaml" && srcPathExt == ".yaml" {
			return fmt.Errorf("destination or source must be binary")
		}

		if dstPathExt != ".yaml" && srcPathExt != ".yaml" {
			return fmt.Errorf("destination or source must be yaml")
		}

		if isSrcArchive {
			fmt.Println("Converting archive", filepath.Base(srcPath), "file", srcFile, "to", filepath.Base(dstPath))
		}
		if isDstArchive {
			fmt.Println("Converting", filepath.Base(srcPath), "to archive", filepath.Base(dstPath), "file", dstFile)
		}
		if !isSrcArchive && !isDstArchive {
			fmt.Println("Converting", filepath.Base(srcPath), "to", filepath.Base(dstPath))
		}
		var srcArchive *pfs.PFS
		var dstArchive *pfs.PFS
		if isSrcArchive {
			srcArchive = &pfs.PFS{}
			srcArchive, err = pfs.NewFile(srcPath)
			if err != nil {
				return fmt.Errorf("open src archive: %w", err)
			}
		}

		if isDstArchive {
			dstArchive = &pfs.PFS{}
			dstArchive, err = pfs.NewFile(srcPath)
			if err != nil {
				return fmt.Errorf("open dst archive: %w", err)
			}
		}

		var srcData []byte

		if srcArchive != nil && srcFile != "" {
			srcData, err = srcArchive.File(srcFile)
			if err != nil {
				return fmt.Errorf("get src file: %w", err)
			}
		}
		if srcArchive == nil {
			srcData, err = os.ReadFile(srcPath)
			if err != nil {
				return fmt.Errorf("open src file: %w", err)
			}
		}

		fileFormat := ""
		isYamlOut := false
		if isSrcArchive || srcPathExt != ".yaml" {
			fileFormat = srcPathExt
			if srcFile != "" {
				fileFormat = filepath.Ext(srcFile)
			}

			isYamlOut = true
		}
		if isDstArchive || dstPathExt != ".yaml" {
			fileFormat = dstPathExt
			if dstFile != "" {
				fileFormat = filepath.Ext(dstFile)
			}
		}

		buf := &bytes.Buffer{}
		if isYamlOut {
			fmt.Println("Source is", len(srcData), "bytes, turning into yaml...")
			reader := raw.New(fileFormat)
			if reader == nil {
				return fmt.Errorf("unsupported file format %s", fileFormat)
			}
			err = reader.Read(bytes.NewReader(srcData))
			if err != nil {
				return fmt.Errorf("read %s: %w", filepath.Base(srcPath), err)
			}

			err = yaml.NewEncoder(buf).Encode(reader)
			if err != nil {
				return fmt.Errorf("yaml encode %s: %w", fileFormat, err)
			}

			w, err := os.Create(dstPath)
			if err != nil {
				return fmt.Errorf("create dst file %s: %w", filepath.Base(dstPath), err)
			}
			defer w.Close()
			if dstArchive != nil && dstFile != "" {
				err = dstArchive.Add(dstFile, buf.Bytes())
				if err != nil {
					return fmt.Errorf("add dst file: %w", err)
				}

				err = dstArchive.Write(w)
				if err != nil {
					return fmt.Errorf("write dst archive: %w", err)
				}
				return nil
			}

			_, err = w.Write(buf.Bytes())
			if err != nil {
				return fmt.Errorf("write %s: %w", filepath.Base(dstPath), err)
			}
			return nil
		}

		writer := raw.New(fileFormat)
		if writer == nil {
			return fmt.Errorf("unsupported file format %s", fileFormat)
		}
		err = yaml.Unmarshal(srcData, writer)
		if err != nil {
			return fmt.Errorf("yaml unmarshal: %w", err)
		}
		err = writer.Write(buf)
		if err != nil {
			return fmt.Errorf("write %s: %w", filepath.Base(dstPath), err)
		}

		log.Infof("Destination is %d bytes, turning into %s...", buf.Len(), fileFormat)
		w, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("create dst file: %w", err)
		}
		defer w.Close()

		if dstArchive != nil && dstFile != "" {
			err = dstArchive.Add(dstFile, buf.Bytes())
			if err != nil {
				return fmt.Errorf("add dst file: %w", err)
			}

			err = dstArchive.Write(w)
			if err != nil {
				return fmt.Errorf("encode dst archive: %w", err)
			}
			return nil
		}

		_, err = w.Write(buf.Bytes())
		if err != nil {
			return fmt.Errorf("write %s: %w", filepath.Base(dstPath), err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(yamlCmd)
	yamlCmd.PersistentFlags().String("path", "", "path to inspect")
	yamlCmd.PersistentFlags().String("file", "", "file to inspect inside pfs")
}
