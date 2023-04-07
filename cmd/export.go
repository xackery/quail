package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/model/mesh/ter"
	"github.com/xackery/quail/model/metadata/zon"
	"github.com/xackery/quail/pfs/archive"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/pfs/s3d"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export an eqg or s3d pfs to quail-addon",
	Long:  `Export an eqg or s3d pfs to quail-addon`,
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

		isBlender, err := cmd.Flags().GetBool("blender")
		if err != nil {
			return fmt.Errorf("parse blender: %w", err)
		}

		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path check: %w", err)
		}
		if fi.IsDir() {
			if isBlender {
				return fmt.Errorf("blender export is not supported for directories")
			}
			if strings.Contains(out, "_") {
				out = fmt.Sprintf("./%s", filepath.Base(path))
			}
			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			pfsCount := 0
			fmt.Println("looking for pfs pfss in", path)
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				fileOut := fmt.Sprintf("%s/_%s", out, file.Name())
				fi, err := os.Stat(fileOut)
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Printf("creating directory %s/\n", fileOut)
						err = os.MkdirAll(fileOut, 0766)
						if err != nil {
							return fmt.Errorf("mkdirall: %w", err)
						}
					}
					fi, err = os.Stat(fileOut)
					if err != nil {
						return fmt.Errorf("stat after mkdirall: %w", err)
					}
				}
				if !fi.IsDir() {
					return fmt.Errorf("%s is not a directory", fileOut)
				}

				inFile := fmt.Sprintf("%s/%s", path, file.Name())

				fmt.Println("parsing", inFile)
				ok, err := export(inFile, fileOut, false)
				if err != nil {
					return err
				}
				if !ok {
					continue
				}
				pfsCount++
			}
			if pfsCount == 0 {
				return fmt.Errorf("no pfss found")
			}
			fmt.Println("exported", pfsCount, "pfss")
			return nil
		}

		fi, err = os.Stat(out)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("creating directory %s/\n", out)
				err = os.MkdirAll(out, 0766)
				if err != nil {
					return fmt.Errorf("mkdirall: %w", err)
				}
			}
			fi, err = os.Stat(out)
			if err != nil {
				return fmt.Errorf("stat after mkdirall: %w", err)
			}
		}
		if !fi.IsDir() {
			return fmt.Errorf("%s is not a directory", out)
		}
		fmt.Println("parsing", path)

		_, err = export(path, out, isBlender)
		if err != nil {
			return err
		}
		return nil
	},
}

func export(in string, out string, isBlender bool) (bool, error) {
	var err error
	switch strings.ToLower(filepath.Ext(in)) {
	case ".eqg":
		err = exportEQG(in, out, isBlender)
		if err != nil {
			return true, fmt.Errorf("exportEQG: %w", err)
		}
	case ".s3d":
		err = exportS3D(in, out, isBlender)
		if err != nil {
			return true, fmt.Errorf("exportS3D: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}

func exportEQG(in string, out string, isBlender bool) error {
	start := time.Now()
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	pfs, err := eqg.New("out")
	if err != nil {
		return fmt.Errorf("eqg.New: %w", err)
	}
	err = pfs.Decode(f)
	if err != nil {
		return fmt.Errorf("decode %s: %w", in, err)
	}
	err = exportBlender(pfs, out)
	if err != nil {
		return fmt.Errorf("exportBlender: %w", err)
	}
	fmt.Printf("%s exported in %.1fs\n", filepath.Base(in), time.Since(start).Seconds())
	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().String("path", "", "path to compressed eqg")
	exportCmd.PersistentFlags().String("out", "", "out folder to export to")
	exportCmd.PersistentFlags().Bool("blender", false, "export to quail-addon format")
}

func exportS3D(path string, out string, isBlender bool) error {
	start := time.Now()
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
	if isBlender {
		err = exportBlender(e, out)
		if err != nil {
			return fmt.Errorf("exportBlender: %w", err)
		}
		fmt.Printf("%s exported in %.1fs\n", filepath.Base(path), time.Since(start).Seconds())
		return nil
	}

	fmt.Printf("%s exported in %.1fs\n", filepath.Base(path), time.Since(start).Seconds())
	return nil
}

func exportBlender(pfs archive.ReadWriter, out string) error {
	var err error
	type modeler interface {
		Decode(r io.ReadSeeker) error
		BlenderExport(dir string) error
	}

	var a modeler
	for _, fe := range pfs.Files() {
		ext := strings.ToLower(filepath.Ext(fe.Name()))
		switch ext {
		case ".ter":
			a, err = ter.New(fe.Name(), pfs)
		case ".zon":
			a, err = zon.New(fe.Name(), pfs)
		case ".png":
			err = os.WriteFile(fmt.Sprintf("%s/%s", out, fe.Name()), fe.Data(), os.ModePerm)
			if err != nil {
				return fmt.Errorf("write %s: %w", fe.Name(), err)
			}
			continue
		default:
			return fmt.Errorf("unsupported file extension %s", ext)
		}
		if err != nil {
			return fmt.Errorf("%s.New %s: %w", ext, fe.Name(), err)
		}
		r := bytes.NewReader(fe.Data())
		err = a.Decode(r)
		if err != nil {
			return fmt.Errorf("%s.Decode %s: %w", ext, fe.Name(), err)
		}
		err = a.BlenderExport(out)
		if err != nil {
			return fmt.Errorf("%s.BlenderExport %s: %w", ext, fe.Name(), err)
		}

	}
	return nil
}
