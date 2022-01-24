/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/blend"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
	"github.com/xackery/quail/zon"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return fmt.Errorf("parse out: %w", err)
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("parse absolute path: %w", err)
		}
		if out == "" {
			if len(args) < 2 {
				out = filepath.Base(absPath)
			} else {
				out = args[1]
			}
		}

		blenderPath, _ := cmd.Flags().GetString("blender")

		err = importExec(path, out, blenderPath)
		if err != nil {
			fmt.Println("import:", err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().String("path", "", "path to import")
	importCmd.PersistentFlags().String("blender", "", "blender path (optional)")
	importCmd.PersistentFlags().String("out", "", "name of compressed eqg archive output, defaults to path's basename")
	importCmd.Example = `quail import --path="./_clz.eqg/"`
}

func importExec(path string, out string, blenderPath string) error {
	start := time.Now()

	out = strings.ToLower(out)

	if strings.Contains(out, ".") && !strings.HasSuffix(out, ".eqg") {
		return fmt.Errorf("only .eqg extension out names are supported")
	}

	if !strings.HasSuffix(out, ".eqg") {
		out = out + ".eqg"
	}
	out = strings.TrimPrefix(out, "_")

	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("path invalid, must be a directory (%s)", path)
	}

	shortname, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("filepath.Abs %s: %w", path, err)
	}
	shortname = strings.TrimSuffix(filepath.Base(shortname), ".eqg")

	fmt.Println("shortname", shortname)
	cachePath := fmt.Sprintf("%s/cache/", path)
	fi, err = os.Stat(cachePath)
	if err != nil {
		err = os.MkdirAll(cachePath, 0766)
		if err != nil {
			fmt.Printf("cache folder not found for import: %s\n", err)
			os.Exit(1)
		}
		fi, err = os.Stat(cachePath)
		if err != nil {
			fmt.Printf("cache folder stat: %s\n", err)
			os.Exit(1)
		}
	}
	if !fi.IsDir() {
		fmt.Println("cache file exists, should be folder")
		os.Exit(1)
	}

	err = blend.Convert(path, shortname, blenderPath)
	if err != nil {
		return fmt.Errorf("blend.Convert: %w", err)
	}
	fi, err = os.Stat(fmt.Sprintf("%s/%s.obj", cachePath, shortname))
	if err != nil {
		return fmt.Errorf("import %s.obj: not found in cache", shortname)
	}
	if fi.IsDir() {
		return fmt.Errorf("import %s.obj: is a directory", shortname)
	}

	ter := &ter.TER{}
	err = ter.ImportObj(fmt.Sprintf("%s/%s.obj", cachePath, shortname), fmt.Sprintf("%s/%s.mtl", cachePath, shortname), fmt.Sprintf("%s/%s_material.txt", cachePath, shortname))
	if err != nil {
		return fmt.Errorf("importObj: %w", err)
	}
	terW, err := os.Create(fmt.Sprintf("%s/%s.ter", cachePath, shortname))
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer terW.Close()
	err = ter.Save(terW)
	if err != nil {
		return fmt.Errorf("ter.Save: %w", err)
	}

	zon := &zon.ZON{}
	err = zon.Import(fmt.Sprintf("%s/%s_light.txt", cachePath, shortname), fmt.Sprintf("%s/%s_mod.txt", cachePath, shortname), fmt.Sprintf("%s/%s_region.txt", cachePath, shortname))
	if err != nil {
		return fmt.Errorf("zon.Import: %w", err)
	}

	zonW, err := os.Create(fmt.Sprintf("%s/%s.zon", cachePath, shortname))
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer zonW.Close()
	err = zon.Save(zonW)
	if err != nil {
		return fmt.Errorf("zon.Save: %w", err)
	}

	modNames := zon.ModelNames()
	for _, modName := range modNames {
		mod := &mod.MOD{}
		modName = strings.TrimSuffix(modName, ".obj")
		err = mod.ImportObj(fmt.Sprintf("%s/%s.obj", cachePath, modName), fmt.Sprintf("%s/%s.mtl", cachePath, modName), fmt.Sprintf("%s/%s_material.txt", cachePath, modName))
		if err != nil {
			return fmt.Errorf("importObj: %w", err)
		}

		modPath := fmt.Sprintf("%s/%s.mod", cachePath, modName)
		modW, err := os.Create(modPath)
		if err != nil {
			return fmt.Errorf("mod create %s: %w", modPath, err)
		}
		err = mod.Save(modW)
		if err != nil {
			return fmt.Errorf("save: %w", err)
		}
		modW.Close()
	}

	err = zon.AddModel(fmt.Sprintf("%s.ter", shortname))
	if err != nil {
		return fmt.Errorf("addModel %s.ter: %w", shortname, err)
	}

	e := &eqg.EQG{}
	files, err := os.ReadDir(cachePath)
	if err != nil {
		return fmt.Errorf("readdir cachepath: %w", err)
	}
	if len(files) == 0 {
		fmt.Printf("no files found in %s to import\n", cachePath)
		os.Exit(1)
	}

	suffixes := []string{".py", ".obj", ".mtl", "_material.txt", "_light.txt", "_region.txt", "_doors.txt", "_mod.txt"}
	addStdout := ""
	fileCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileCount++
		isIgnored := false
		for _, suffix := range suffixes {
			if !strings.HasSuffix(file.Name(), suffix) {
				continue
			}
			isIgnored = true
			break
		}
		if isIgnored {
			continue
		}

		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", cachePath, file.Name()))
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name(), err)
		}
		err = e.Add(file.Name(), data)
		if err != nil {
			return fmt.Errorf("add %s: %w", file.Name(), err)
		}
		addStdout += file.Name() + ", "
	}
	if fileCount == 0 {
		fmt.Println("no files found to add")
		os.Exit(1)
	}
	addStdout = addStdout[0:len(addStdout)-2] + "\n"

	w, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("create %s: %w", out, err)
	}
	defer w.Close()
	err = e.Save(w)
	if err != nil {
		return fmt.Errorf("save %s: %w", out, err)
	}

	fmt.Printf("%s created with %d file%s: %s", out, fileCount, helper.Pluralize(fileCount), addStdout)
	fmt.Printf("import took %0.2f seconds\n", time.Since(start).Seconds())
	return nil
}
