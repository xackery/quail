package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
	"github.com/xackery/quail/model/metadata/ani"
	"github.com/xackery/quail/model/metadata/zon"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/pfs/s3d"
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

		filter, err := cmd.Flags().GetString("filter")
		if err != nil {
			return fmt.Errorf("parse filter: %w", err)
		}
		if filter == "" {
			if len(args) < 2 {
				filter = "all"
			} else {
				filter = args[1]
			}
		}

		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}

		if out == "" {
			if len(args) < 3 {
				out = fmt.Sprintf("debug_%s", filepath.Base(path))
			} else {
				out = args[2]
			}
		}
		out = strings.TrimSuffix(out, ".png")
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
			{instance: &ani.ANI{}, extension: ".ani"},
			{instance: &eqg.EQG{}, extension: ".eqg"},
			{instance: &s3d.S3D{}, extension: ".s3d"},
			{instance: &mod.MOD{}, extension: ".mod"},
			{instance: &mds.MDS{}, extension: ".mds"},
			{instance: &ter.TER{}, extension: ".ter"},
			{instance: &zon.ZON{}, extension: ".zon"},
		}

		fmt.Println("debugging file", path, "and dumping results to", out)
		if filter != "all" {
			fmt.Println("filtering by", filter)
		}
		for _, v := range decodes {
			if ext != v.extension {
				continue
			}

			if ext == ".eqg" {
				err = debugEQG(path, out, filter)
				if err != nil {
					return fmt.Errorf("debugEQG: %w", err)
				}
			}

			if filter != "all" && !strings.Contains(path, filter) {
				continue
			}

			err = dumpDecode(f, v.extension, path, fmt.Sprintf("%s%s", out, filepath.Base(path)))
			if err != nil {
				return fmt.Errorf("dumpDecode: %w", err)
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

func debugEQG(path string, out string, filter string) error {
	archive, err := eqg.New(filepath.Base(path))
	if err != nil {
		return fmt.Errorf("new: %w", err)
	}
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()
	err = archive.Decode(r)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	fmt.Printf("%s contains %d files:\n", filepath.Base(path), archive.Len())

	filesByName := archive.Files()
	sort.Sort(common.FilerByName(filesByName))
	for i, fe := range archive.Files() {
		if filter != "all" && !strings.Contains(fe.Name(), filter) {
			continue
		}
		base := float64(len(fe.Data()))
		strSize := ""
		num := float64(1024)
		if base < num*num*num*num {
			strSize = fmt.Sprintf("%0.0fG", base/num/num/num)
		}
		if base < num*num*num {
			strSize = fmt.Sprintf("%0.0fM", base/num/num)
		}
		if base < num*num {
			strSize = fmt.Sprintf("%0.0fK", base/num)
		}
		if base < num {
			strSize = fmt.Sprintf("%0.0fB", base)
		}

		fmt.Printf("%d: %s %s\n", i, fe.Name(), strSize)
		ext := strings.ToLower(filepath.Ext(fe.Name()))
		r := bytes.NewReader(fe.Data())
		err = dumpDecode(r, ext, fe.Name(), fmt.Sprintf("%s.%s", out, fe.Name()))
		if err != nil {
			return fmt.Errorf("dumpDecode %s: %w", fe.Name(), err)
		}
		fmt.Printf("%s\t%s\n", out, fe.Name())
	}

	return nil
}

func dumpDecode(r io.ReadSeeker, ext string, path string, out string) error {
	var err error
	type decoder interface {
		Decode(io.ReadSeeker) error
	}
	type decodeTypes struct {
		instance  decoder
		extension string
	}
	decodes := []*decodeTypes{
		{instance: &ani.ANI{}, extension: ".ani"},
		{instance: &eqg.EQG{}, extension: ".eqg"},
		{instance: &s3d.S3D{}, extension: ".s3d"},
		{instance: &mod.MOD{}, extension: ".mod"},
		{instance: &mds.MDS{}, extension: ".mds"},
		{instance: &ter.TER{}, extension: ".ter"},
		{instance: &zon.ZON{}, extension: ".zon"},
	}

	for _, v := range decodes {
		if ext != v.extension {
			continue
		}

		fmt.Println("dumping", path)
		dump.New(path)
		err = v.instance.Decode(r)
		if err != nil {
			return fmt.Errorf("failed to decode %s: %w", v.extension, err)
		}
		dump.WriteFileClose(fmt.Sprintf("%s.png", out))
		return nil
	}
	return nil
}
