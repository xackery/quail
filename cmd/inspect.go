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
	"github.com/xackery/quail/lit"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/s3d"
	"github.com/xackery/quail/ter"
	"github.com/xackery/quail/zon"
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
		file, err := cmd.Flags().GetString("file")
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

		var archive common.ArchiveReadWriter
		ext := filepath.Ext(path)
		switch ext {
		case ".eqg":
			e, err := eqg.New(filepath.Base(path))
			if err != nil {
				return fmt.Errorf("eqg new: %w", err)
			}

			if file == "" {
				err = inspectEQG(path)
				if err != nil {
					return fmt.Errorf("inspectEQG: %w", err)
				}
				os.Exit(0)
			}
			r, err := os.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()
			err = e.Decode(r)
			if err != nil {
				return fmt.Errorf("decode: %w", err)
			}

			archive = e
		case ".s3d":
			e, err := s3d.New(filepath.Base(path))
			if err != nil {
				return fmt.Errorf("s3d new: %w", err)
			}
			if file == "" {
				err = inspectS3D(path)
				if err != nil {
					return fmt.Errorf("inspectS3D: %w", err)
				}
				os.Exit(0)
			}
			r, err := os.Open(path)
			if err != nil {
				return err
			}
			defer r.Close()
			err = e.Decode(r)
			if err != nil {
				return fmt.Errorf("decode: %w", err)
			}
			archive = e
		default:
			file = filepath.Base(path)
			archive, err = common.NewPath(path)
			if err != nil {
				return fmt.Errorf("path new: %w", err)
			}
		}

		err = inspect(archive, file)
		if err != nil {
			return fmt.Errorf("inspect: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.PersistentFlags().String("path", "", "path to inspect")
	inspectCmd.PersistentFlags().String("file", "", "file to inspect inside archive")
}

func inspect(archive common.ArchiveReadWriter, file string) error {

	var err error
	ext := strings.ToLower(filepath.Ext(file))

	callbacks := []struct {
		invoke func(file string, archive common.ArchiveReadWriter) error
		name   string
	}{
		{invoke: inspectMDS, name: "mds"},
		{invoke: inspectZON, name: "zon"},
		{invoke: inspectMOD, name: "mod"},
		{invoke: inspectTER, name: "ter"},
		{invoke: inspectLIT, name: "lit"},
	}

	for _, evt := range callbacks {
		if ext != "."+evt.name {
			continue
		}
		err = evt.invoke(file, archive)
		if err != nil {
			return fmt.Errorf("inspect %s: %w", evt.name, err)
		}
		os.Exit(0)
	}

	return fmt.Errorf("unsupported extension: %s", ext)
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
	err = archive.Decode(r)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	fmt.Printf("%s contains %d files:\n", filepath.Base(path), archive.Len())

	filesByName := archive.Files()
	noteworthyFile := "file.ext"

	sort.Sort(common.FilerByName(filesByName))
	for _, fe := range archive.Files() {
		base := float64(len(fe.Data()))
		out := ""
		num := float64(1024)
		if strings.HasSuffix(fe.Name(), ".zon") {
			noteworthyFile = fe.Name()
		}
		if strings.HasSuffix(fe.Name(), ".mds") && noteworthyFile == "file.ext" {
			noteworthyFile = fe.Name()
		}
		if strings.HasSuffix(fe.Name(), ".mod") && noteworthyFile == "file.ext" {
			noteworthyFile = fe.Name()
		}
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

	fmt.Printf("you can inspect files, e.g.: quail inspect %s %s\n", path, noteworthyFile)
	return nil
}

func inspectS3D(path string) error {
	archive, err := s3d.New(filepath.Base(path))
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

func inspectMDS(file string, archive common.ArchiveReadWriter) error {
	e, err := mds.NewFile(filepath.Base(file), archive, file)
	if err != nil {
		return fmt.Errorf("mds new: %w", err)
	}

	e.Inspect()
	return nil
}

func inspectZON(file string, archive common.ArchiveReadWriter) error {
	e, err := zon.NewFile(filepath.Base(file), archive, file)
	if err != nil {
		return fmt.Errorf("zon new: %w", err)
	}

	e.Inspect()
	return nil
}

func inspectMOD(file string, archive common.ArchiveReadWriter) error {
	e, err := mod.NewFile(filepath.Base(file), archive, file)
	if err != nil {
		return fmt.Errorf("mod new: %w", err)
	}

	e.Inspect()
	return nil
}

func inspectTER(file string, archive common.ArchiveReadWriter) error {
	e, err := ter.NewFile(filepath.Base(file), archive, file)
	if err != nil {
		return fmt.Errorf("ter new: %w", err)
	}
	e.Inspect()

	return nil
}

func inspectLIT(file string, archive common.ArchiveReadWriter) error {
	e, err := lit.NewFile(filepath.Base(file), archive, file)
	if err != nil {
		return fmt.Errorf("lit new: %w", err)
	}
	e.Inspect()

	return nil
}
