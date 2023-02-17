package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	qexport "github.com/xackery/quail/model/plugin/export"
	"github.com/xackery/quail/model/plugin/gltf"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/pfs/s3d"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export an eqg or s3d archive to embedded GLTF",
	Long:  `Export an eqg or s3d archive to embedded GLTF`,
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
			fmt.Println("looking for pfs archives in", path)
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
				ok, err := export(inFile, fileOut)
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
			fmt.Println("exported", archiveCount, "archives")
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
		_, err = export(path, out)
		if err != nil {
			return err
		}
		return nil
	},
}

func export(in string, out string) (bool, error) {

	var err error
	switch strings.ToLower(filepath.Ext(in)) {
	case ".eqg":
		err = exportEQG(in, out)
		if err != nil {
			return true, fmt.Errorf("exportEQG: %w", err)
		}
	case ".s3d":
		err = exportS3D(in, out)
		if err != nil {
			return true, fmt.Errorf("exportS3D: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}

func exportEQG(in string, out string) error {
	start := time.Now()
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	archive, err := eqg.New("out")
	if err != nil {
		return fmt.Errorf("eqg.New: %w", err)
	}
	err = archive.Decode(f)
	if err != nil {
		return fmt.Errorf("decode %s: %w", in, err)
	}

	e, err := qexport.New(strings.TrimSuffix(filepath.Base(in), ".eqg"), archive)
	if err != nil {
		return fmt.Errorf("export new: %w", err)
	}

	err = e.LoadArchive()
	if err != nil {
		return fmt.Errorf("decode archive: %w", err)
	}

	outFile := fmt.Sprintf("%s/%s.gltf", out, e.Name())
	fmt.Println("exporting", outFile)
	w, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer w.Close()
	doc, err := gltf.New()
	if err != nil {
		return fmt.Errorf("gltf.New: %w", err)
	}
	err = e.GLTFEncode(doc)
	if err != nil {
		return fmt.Errorf("gltf: %w", err)
	}
	err = doc.Export(w)
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	fmt.Printf("%s exported in %.1fs\n", filepath.Base(in), time.Since(start).Seconds())
	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().String("path", "", "path to compressed eqg")
	exportCmd.PersistentFlags().String("out", "", "out folder to export to")
}

func exportS3D(path string, out string) error {
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

	fmt.Printf("%s exported in %.1fs\n", filepath.Base(path), time.Since(start).Seconds())
	return nil
}
