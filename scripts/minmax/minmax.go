package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
	"github.com/xackery/quail/pfs/eqg"

	"net/http"

	_ "net/http/pprof"
)

var (
	properties = make(map[string]*propertyMinMax)
	propValues = make(map[string]map[string]string)
)

type propertyMinMax struct {
	name string
	min  float32
	max  float32
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("127.0.0.1:8082", nil))
	}()
	start := time.Now()
	err := run()
	fmt.Printf("Finished in %0.2f seconds", time.Since(start).Seconds())
	if err != nil {
		fmt.Println("Failed to run:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		fmt.Println("usage: minmax <path>")
		os.Exit(1)
	}
	path := os.Args[1]
	fmt.Println("Path:", path)
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", path, err)
	}

	for _, eqgFile := range files {
		if eqgFile.IsDir() {
			continue
		}
		if !strings.HasSuffix(eqgFile.Name(), ".eqg") {
			continue
		}
		fmt.Println(eqgFile.Name())
		pfs, err := eqg.NewFile(path + "/" + eqgFile.Name())
		if err != nil {
			return fmt.Errorf("new file %s: %w", eqgFile.Name(), err)
		}
		for _, modFile := range pfs.Files() {
			ext := filepath.Ext(modFile.Name())
			switch ext {
			case ".mod":
				fmt.Println(modFile.Name())
				mesh, err := mod.NewFile(modFile.Name(), pfs, modFile.Name())
				if err != nil {
					fmt.Printf("Failed new mod %s: %s\n", modFile.Name(), err)
					continue
					//return fmt.Errorf("new mod %s: %w", modFile.Name(), err)
				}
				parseMaterials(eqgFile.Name(), modFile.Name(), mesh.MaterialManager.Materials())
				mesh.Close()
			case ".mds":
				fmt.Println(modFile.Name())
				mesh, err := mds.NewFile(modFile.Name(), pfs, modFile.Name())
				if err != nil {
					fmt.Printf("Failed new mod %s: %s\n", modFile.Name(), err)
					continue
					//return fmt.Errorf("new mod %s: %w", modFile.Name(), err)
				}
				parseMaterials(eqgFile.Name(), modFile.Name(), mesh.MaterialManager.Materials())
				mesh.Close()
			case ".ter":
				fmt.Println(modFile.Name())
				mesh, err := ter.NewFile(modFile.Name(), pfs, modFile.Name())
				if err != nil {
					fmt.Printf("Failed new mod %s: %s\n", modFile.Name(), err)
					continue
					//return fmt.Errorf("new mod %s: %w", modFile.Name(), err)
				}
				parseMaterials(eqgFile.Name(), modFile.Name(), mesh.MaterialManager.Materials())
				mesh.Close()
			}
		}

	}
	w, err := os.Create("minmax.txt")
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer w.Close()
	for _, prop := range properties {
		fmt.Fprintf(w, "%s: %0.2f - %0.2f\n", prop.name, prop.min, prop.max)
	}
	for name, values := range propValues {
		fmt.Fprintf(w, "%s:\n", name)
		for k, value := range values {
			fmt.Fprintf(w, "\t%s\t%s\n", k, value)
		}
	}

	return nil
}

func parseMaterials(eqgName string, meshName string, materials []*geo.Material) {
	for _, material := range materials {
		for _, prop := range material.Properties {
			pl, ok := propValues[fmt.Sprintf("%s|%d", prop.Name, prop.Category)]
			if !ok {
				pl = make(map[string]string)
			}
			pl[prop.Value] = fmt.Sprintf("%s|%s|%s", material.ShaderName, eqgName, meshName)
			propValues[fmt.Sprintf("%s|%d", prop.Name, prop.Category)] = pl
			if prop.Category != 0 {
				continue
			}
			value := helper.AtoF32(prop.Value)
			pm, ok := properties[prop.Name]
			if !ok {
				pm = &propertyMinMax{
					name: prop.Name,
					min:  value,
					max:  value,
				}
				properties[prop.Name] = pm
			}
			if value < pm.min {
				pm.min = value
			}
			if value > pm.max {
				pm.max = value
			}

		}
	}
}
