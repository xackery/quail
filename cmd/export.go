package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
	"github.com/xackery/quail/lay"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/s3d"
	"github.com/xackery/quail/zon"
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
				//gltfOut := fmt.Sprintf("%s/_%s/%s", path, file.Name(), strings.TrimPrefix(strings.TrimSuffix(strings.TrimSuffix(filepath.Base(fileOut), ".eqg"), ".s3d"), "_")+".gltf")

				//fmt.Println("exporting", inFile, "to", gltfOut)
				fmt.Println("parsing", inFile)
				ok, err := export(inFile, fileOut, true)
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
		_, err = export(path, out, false)
		if err != nil {
			return err
		}
		return nil
	},
}

func export(in string, out string, isDir bool) (bool, error) {

	var err error
	switch strings.ToLower(filepath.Ext(in)) {
	case ".eqg":
		err = exportEQG(in, out, isDir)
		if err != nil {
			return true, fmt.Errorf("exportEQG: %w", err)
		}
	case ".s3d":
		err = exportS3D(in, out, isDir)
		if err != nil {
			return true, fmt.Errorf("exportS3D: %w", err)
		}
	default:
		return false, nil
	}
	return true, nil
}

func exportEQG(in string, out string, isDir bool) error {
	start := time.Now()
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	a, err := eqg.New("out")
	if err != nil {
		return fmt.Errorf("eqg.New: %w", err)
	}
	err = a.Load(f)
	if err != nil {
		return fmt.Errorf("load %s: %w", in, err)
	}

	isZone := true
	zoneName := fmt.Sprintf("%s.zon", strings.TrimSuffix(filepath.Base(in), ".eqg"))
	zoneData, err := a.File(zoneName)
	if err != nil {
		isZone = false
	}

	if !isZone {
		for _, fe := range a.Files() {
			err = convertFile(a, fe.Name(), fe.Data(), out)
			if err != nil {
				return fmt.Errorf("convert %s: %w", fe.Name(), err)
			}
		}
	} else {
		err = convertFile(a, zoneName, zoneData, out)
		if err != nil {
			return fmt.Errorf("convert %s: %w", zoneName, err)
		}
	}

	fmt.Printf("%s exported in %.1fs\n", filepath.Base(in), time.Since(start).Seconds())
	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().String("path", "", "path to compressed eqg")
	exportCmd.PersistentFlags().String("out", "", "out folder to export to")
}

func exportS3D(path string, out string, isDir bool) error {
	start := time.Now()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	e, err := s3d.New("out")
	if err != nil {
		return fmt.Errorf("eqg.New: %w", err)
	}
	err = e.Load(f)
	if err != nil {
		return fmt.Errorf("load %s: %w", path, err)
	}

	fmt.Printf("%s exported in %.1fs\n", filepath.Base(path), time.Since(start).Seconds())
	return nil
}

func convertFile(archive common.Archiver, name string, data []byte, out string) error {

	outFile := ""
	switch strings.ToLower(filepath.Ext(name)) {
	case ".mds":
		e, err := mds.NewEQG(name, archive)
		if err != nil {
			return fmt.Errorf("newEQG: %w", err)
		}

		err = e.Load(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("load: %w", err)
		}

		err = layerInject(archive, e, fmt.Sprintf("%s.lay", strings.TrimSuffix(name, ".mds")))
		if err != nil {
			return fmt.Errorf("layerInject: %w", err)
		}

		outFile = fmt.Sprintf("%s/%s.gltf", out, strings.TrimSuffix(name, ".mds"))
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
		err = e.GLTFExport(doc)
		if err != nil {
			return fmt.Errorf("gltf: %w", err)
		}
		err = doc.Export(w)
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}

	case ".mod":
		e, err := mod.NewEQG(name, archive)
		if err != nil {
			return fmt.Errorf("newEQG: %w", err)
		}

		err = e.Load(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("load: %w", err)
		}

		err = layerInject(archive, e, fmt.Sprintf("%s.lay", strings.TrimSuffix(name, ".mod")))
		if err != nil {
			return fmt.Errorf("layerInject: %w", err)
		}
		outFile = fmt.Sprintf("%s/%s.gltf", out, strings.TrimSuffix(name, ".mod"))
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
		err = e.GLTFExport(doc)
		if err != nil {
			return fmt.Errorf("gltf: %w", err)
		}
		err = doc.Export(w)
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}

	case ".ter":
		// we skip terrain data in archive, and instead load it via .zon
	case ".zon":
		z, err := zon.NewEQG(name, archive)
		if err != nil {
			return fmt.Errorf("new: %w", err)
		}

		err = z.Load(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("load: %w", err)
		}

		doc, err := gltf.New()
		if err != nil {
			return fmt.Errorf("gltf.New: %w", err)
		}
		err = z.GLTFExport(doc)
		if err != nil {
			return fmt.Errorf("glts: %w", err)
		}
		outFile = fmt.Sprintf("%s/%s.gltf", out, strings.TrimSuffix(name, ".zon"))
		fmt.Println("exporting", outFile)
		w, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}
		defer w.Close()
		err = doc.Export(w)
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}

	}
	return nil
}

func layerInject(archive common.Archiver, modeler common.Modeler, layName string) error {

	layEntry, err := archive.File(layName)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil
		}
		return fmt.Errorf("file %s: %w", layName, err)
	}

	if len(layEntry) == 0 {
		return nil
	}

	l, err := lay.NewEQG(layName, archive)
	if err != nil {
		return fmt.Errorf("lay.NewEQG: %w", err)
	}
	err = l.Load(bytes.NewReader(layEntry))
	if err != nil {
		return fmt.Errorf("lay.Load: %w", err)
	}
	err = modeler.SetLayers(l.Layers())
	if err != nil {
		return fmt.Errorf("setLayers: %w", err)
	}
	return nil
}
