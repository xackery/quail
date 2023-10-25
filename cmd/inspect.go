package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/metadata/ani"
	"github.com/xackery/quail/model/metadata/lay"
	"github.com/xackery/quail/model/metadata/prt"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/zon"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/quail"
	"gopkg.in/yaml.v3"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect an EverQuest asset",
	Long: `Inspect an EverQuest asset to discover contents within
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

		inspected, err := inspect(path, file)
		if err != nil {
			return fmt.Errorf("inspect: %w", err)
		}

		buf := &bytes.Buffer{}
		enc := yaml.NewEncoder(buf)
		enc.SetIndent(2)
		err = enc.Encode(inspected)
		if err != nil {
			return fmt.Errorf("yaml encode: %w", err)
		}

		log.Infoln(buf.String())

		//reflectTraversal(inspected, 0, -1)
		return nil

	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.PersistentFlags().String("path", "", "path to inspect")
	inspectCmd.PersistentFlags().String("file", "", "file to inspect inside pfs")
}

func reflectTraversal(inspected interface{}, nest int, index int) {
	v := reflect.ValueOf(inspected)
	tv := v.Type()

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		tv = v.Type()
	}

	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			log.Infof("%s%s (Empty)", strings.Repeat("  ", nest), tv.Name())
			return
		}
		log.Infof("%s%s:", strings.Repeat("  ", nest), tv.Name())
		for i := 0; i < v.Len(); i++ {
			reflectTraversal(v.Index(i).Interface(), nest+1, i)
		}
		return
	}

	if v.Kind() != reflect.Struct {
		log.Infof("%s%v", strings.Repeat("  ", nest), v.Interface())
		return
	}

	for i := 0; i < v.NumField(); i++ {
		// check if field is exported
		if tv.Field(i).PkgPath != "" {
			continue
		}

		indexStr := ""
		if index >= 0 {
			indexStr = fmt.Sprintf("[%d]", index)
		}

		// is it a slice?
		if v.Field(i).Kind() == reflect.Slice {
			s := v.Field(i)
			if s.Len() == 0 {
				log.Infof("%s%s %s: (Empty)", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name)
				continue
			}
			log.Infof("%s%s %s:", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name)

			for j := 0; j < s.Len(); j++ {
				if tv.Field(i).PkgPath != "" {
					continue
				}

				if s.Index(j).Kind() == reflect.Uint8 {
					if j == 0 {
						fmt.Printf("%s", strings.Repeat("  ", nest+1))
					}
					fmt.Printf("0x%02x ", s.Index(j).Interface())
					if j == s.Len()-1 {
						fmt.Println()
					}
					continue
				}
				reflectTraversal(s.Index(j).Interface(), nest+1, j)
				//log.Infof("  %d %s\t %+v", j, tv.Field(i).Name, s.Index(j).Interface())
			}
			continue
		}

		if tv.Field(i).Name == "MaterialName" {
			continue
		}

		log.Infof("%s%s %s: %v", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name, v.Field(i).Interface())
	}
}

func inspect(path string, file string) (interface{}, error) {
	/* if len(file) < 2 {
		//log.Infof("Inspecting %s", filepath.Base(path))
	} else {
		//log.Infof("Inspecting %s %s", filepath.Base(path), filepath.Base(file))
	} */

	isValidExt := false
	exts := []string{".eqg", ".s3d", ".pfs", ".pak"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return inspectFile(nil, path, file)
	}

	pfs, err := pfs.NewFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s load: %w", ext, err)
	}
	if len(file) < 2 {
		pfs.Close()
		return *pfs, nil
	}
	return inspectFile(pfs, path, file)
}

func inspectFile(pfs *pfs.PFS, path string, file string) (interface{}, error) {
	if pfs == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return inspectContent(filepath.Base(file), bytes.NewReader(data))
	}

	for _, fe := range pfs.Files() {
		if !strings.EqualFold(fe.Name(), file) {
			continue
		}
		return inspectContent(file, bytes.NewReader(fe.Data()))
	}
	return nil, fmt.Errorf("%s not found in %s", file, filepath.Base(path))
}

func inspectContent(file string, data *bytes.Reader) (interface{}, error) {
	var err error
	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".ani":
		anim := common.NewAnimation(strings.TrimSuffix(strings.ToUpper(file), ".ANI"))
		err = ani.Decode(anim, data)
		if err != nil {
			return nil, fmt.Errorf("ani.Decode %s: %w", file, err)
		}
		return anim, nil
	case ".mds":
		model := common.NewModel(strings.TrimSuffix(strings.ToUpper(file), ".MDS"))
		err = mds.Decode(model, data)
		if err != nil {
			return nil, fmt.Errorf("mds.Decode %s: %w", file, err)
		}
		return model, nil
	case ".mod":
		model := common.NewModel(strings.TrimSuffix(strings.ToUpper(file), ".MOD"))
		err = mod.Decode(model, data)
		if err != nil {
			return nil, fmt.Errorf("mod.Decode %s: %w", file, err)
		}
		return model, nil
	case ".pts":
		point := common.NewParticlePoint(strings.TrimSuffix(strings.ToUpper(file), ".MDS"))
		err = pts.Decode(point, data)
		if err != nil {
			return nil, fmt.Errorf("pts.Decode %s: %w", file, err)
		}
		return point, nil
	case ".prt":
		render := common.NewParticleRender(strings.TrimSuffix(strings.ToUpper(file), ".MDS"))
		err = prt.Decode(render, data)
		if err != nil {
			return nil, fmt.Errorf("prt.Decode %s: %w", file, err)
		}
		return render, nil
	case ".zon":
		zone := common.NewZone(strings.TrimSuffix(strings.ToUpper(file), ".ZON"))
		err = zon.Decode(zone, data)
		if err != nil {
			return nil, fmt.Errorf("zon.Decode %s: %w", file, err)
		}
		return zone, nil
	case ".wld":
		q := quail.New()
		world, err := q.WLDDecode(data, nil)
		if err != nil {
			return nil, fmt.Errorf("wld.Decode %s: %w", file, err)
		}
		return world, nil
	case ".lay":
		model := common.NewModel("")
		err := lay.Decode(model, data)
		if err != nil {
			return nil, fmt.Errorf("lay.Decode %s: %w", file, err)
		}
		return model.Layers, nil
	default:
		return nil, fmt.Errorf("unknown file type %s", ext)
	}
}
