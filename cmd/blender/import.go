package blender

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
	"github.com/xackery/quail/model/mesh/wld"
	"github.com/xackery/quail/model/metadata/ani"
	"github.com/xackery/quail/model/metadata/lit"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/tog"
	"github.com/xackery/quail/pfs/archive"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/pfs/s3d"
)

// blenderImportCmd represents the blenderImport command
var BlenderImportCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		defer func() {
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}()
		err = import_blender(cmd, args)
		return err
	},
}

func init() {
	BlenderImportCmd.PersistentFlags().String("out", "", "name of compressed eqg archive output, defaults to path's basename")
}

func import_blender(cmd *cobra.Command, args []string) error {
	var err error
	start := time.Now()
	if err != nil {
		return fmt.Errorf("parse path: %w", err)
	}
	var path string
	if path == "" {
		if len(args) < 1 {
			return cmd.Usage()
		}
		path = args[0]
	}
	//out, err := cmd.Flags().GetString("out")
	//if err != nil {
	//	return fmt.Errorf("parse out: %w", err)
	//}

	out := ""
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

	out = strings.ToLower(out)

	if strings.Contains(out, ".") && !strings.HasSuffix(out, ".eqg") && !strings.HasSuffix(out, ".s3d") {
		return fmt.Errorf("only eqg and s3d extension out names are supported")
	}

	out = strings.TrimPrefix(out, "_")

	var pfs archive.ReadWriter

	switch strings.ToLower(filepath.Ext(path)) {
	case ".eqg":
		pfs, err = eqg.New(out)
		if err != nil {
			return fmt.Errorf("eqg.New: %w", err)
		}
	case ".s3d":
		pfs, err = s3d.New(out)
		if err != nil {
			return fmt.Errorf("s3d.New: %w", err)
		}
	default:
		return fmt.Errorf("unknown file type %s, only .eqg and .s3d supported", filepath.Ext(path))
	}

	fmt.Println(path, out)
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("os.ReadDir: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found")
	}

	type modeler interface {
		Encode(io.Writer) error
		BlenderImport(out string) error
	}

	for _, fe := range files {
		var e modeler
		var data []byte

		name := fe.Name()
		name = strings.TrimPrefix(name, "_")

		ext := strings.ToLower(filepath.Ext(fe.Name()))

		isRawFile := false
		switch ext {
		case ".mds":
			e, err = mds.New(name, pfs)
		case ".mod":
			e, err = mod.New(name, pfs)
		case ".ani":
			e, err = ani.New(name, pfs)
		case ".pts":
			e, err = pts.New(name, pfs)
		case ".lit":
			e, err = lit.New(name, pfs)
		case ".tog":
			e, err = tog.New(name, pfs)
		case ".ter":
			e, err = ter.New(name, pfs)
		case ".wld":
			e, err = wld.New(name, pfs)
		//case ".lod":
		//	e, err = lod.New(name, pfs)
		case ".prt":
			//e, err = pts.New(name, pfs)
			fmt.Println("TODO: prt PTCL support")
			continue
		case ".dds":
			isRawFile = true
		case ".bmp":
			isRawFile = true
		case ".png":
			isRawFile = true
		case ".txt":
			isRawFile = true
		default:
			fmt.Println("TODO:", fe.Name(), "support")
			continue
			//return fmt.Errorf("unknown file type %s", fe.Name())
		}
		if isRawFile {
			data, err = os.ReadFile(path + "/" + fe.Name())
			if err != nil {
				return fmt.Errorf("dds.ReadFile %s: %w", fe.Name(), err)
			}
			err = pfs.WriteFile(name, data)
			if err != nil {
				return fmt.Errorf("dds.WriteFile %s: %w", fe.Name(), err)
			}
			continue
		}

		if err != nil {
			return fmt.Errorf("%s.New %s: %w", ext, fe.Name(), err)
		}
		err = e.BlenderImport(path + "/" + fe.Name())
		if err != nil {
			return fmt.Errorf("%s.BlenderImport %s: %w", ext, fe.Name(), err)
		}
		buf := &bytes.Buffer{}
		err = e.Encode(buf)
		if err != nil {
			return fmt.Errorf("%s.Encode %s: %w", ext, fe.Name(), err)
		}
		err = pfs.WriteFile(name, buf.Bytes())
		if err != nil {
			return fmt.Errorf("%s.WriteFile %s: %w", ext, fe.Name(), err)
		}
		fmt.Println(name)
	}

	w, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	err = pfs.Encode(w)
	if err != nil {
		return fmt.Errorf("pfs.Encode: %w", err)
	}

	fmt.Printf("%s exported in %.1fs\n", filepath.Base(out), time.Since(start).Seconds())
	return nil
}
