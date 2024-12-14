package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/os"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func init() {
	rootCmd.AddCommand(ondemandCmd)
	ondemandCmd.PersistentFlags().String("path", "", "path to ondemand")
	ondemandCmd.PersistentFlags().String("file", "", "file to ondemand inside pfs")
}

// ondemandCmd represents the ondemand command
var ondemandCmd = &cobra.Command{
	Use:   "ondemand",
	Short: "Generate OnDemandResources.txt entries",
	Long:  `Run against a eqg and generate OnDemandResources.txt`,
	Run:   runOndemand,
}

func runOndemand(cmd *cobra.Command, args []string) {
	err := runOndemandE(cmd, args)
	if err != nil {
		fmt.Printf("Failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func runOndemandE(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("ondemand requires a target file, directory provided")
	}

	err = ondemand(path)
	if err != nil {
		return fmt.Errorf("ondemand: %w", err)
	}

	//reflectTraversal(ondemanded, 0, -1)
	return nil

}

func ondemand(path string) error {
	/* if len(file) < 2 {
		//log.Infof("Ondemanding %s", filepath.Base(path))
	} else {
		//log.Infof("Ondemanding %s %s", filepath.Base(path), filepath.Base(file))
	} */

	isValidExt := false
	exts := []string{".eqg"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return fmt.Errorf("invalid extension %s", ext)
	}

	pfs, err := pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("%s load: %w", ext, err)
	}

	eqgName := strings.ToLower(filepath.Base(path))
	for _, fe := range pfs.Files() {
		feName := strings.ToUpper(fe.Name())
		ext := strings.ToLower(filepath.Ext(fe.Name()))
		switch ext {
		case ".ani":
			fmt.Printf("%s^%s^%s^EQGA\n", eqgName, feName, strings.TrimSuffix(feName, filepath.Ext(feName)))
		case ".lay":
			lay := &raw.Lay{}
			err = lay.Read(bytes.NewReader(fe.Data()))
			if err != nil {
				return fmt.Errorf("lay read: %w", err)
			}

			for _, le := range lay.Entries {
				fmt.Printf("%s^%s^%s^EQGL\n", eqgName, feName, le.Material)
			}
		case ".mod":
			fmt.Printf("%s^%s^%s_ACTORDEF^EQGM\n", eqgName, feName, strings.TrimSuffix(feName, filepath.Ext(feName)))
		case ".mds":
			fmt.Printf("%s^%s^%s_ACTORDEF^EQGS\n", eqgName, feName, strings.TrimSuffix(feName, filepath.Ext(feName)))
		default:
			continue
		}
	}

	return nil
}
