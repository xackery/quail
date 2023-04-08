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
	"github.com/xackery/quail/model/metadata/ani"
	"github.com/xackery/quail/model/metadata/lit"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/tog"
	"github.com/xackery/quail/pfs/archive"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/pfs/s3d"
)

// blenderExportCmd represents the blenderExport command
var BlenderExportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return export(cmd, args)
	},
}

func export(cmd *cobra.Command, args []string) error {
	start := time.Now()
	if len(args) < 1 {
		return cmd.Usage()
	}
	path := args[0]
	out := "./"
	if len(args) > 1 {
		out = args[1]
	}

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("path must be a file, not a directory")
	}

	var pfs archive.ReadWriter

	out += "/_" + filepath.Base(path)
	err = os.MkdirAll(out, 0755)
	if err != nil {
		return fmt.Errorf("os.MkdirAll %s: %w", out, err)
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".eqg":
		pfs, err = eqg.NewFile(path)
		if err != nil {
			return fmt.Errorf("eqg.New: %w", err)
		}
	case ".s3d":
		pfs, err = s3d.NewFile(path)
		if err != nil {
			return fmt.Errorf("s3d.New: %w", err)
		}
	default:
		return fmt.Errorf("unknown file type %s, only .eqg and .s3d supported", filepath.Ext(path))
	}

	type modeler interface {
		Decode(io.ReadSeeker) error
		BlenderExport(out string) error
	}

	for _, fe := range pfs.Files() {
		var e modeler
		ext := strings.ToLower(filepath.Ext(fe.Name()))
		switch ext {
		case ".mds":
			e, err = mds.New(fe.Name(), pfs)
		case ".mod":
			e, err = mod.New(fe.Name(), pfs)
		case ".ani":
			e, err = ani.New(fe.Name(), pfs)
		case ".pts":
			e, err = pts.New(fe.Name(), pfs)
		case ".lit":
			e, err = lit.New(fe.Name(), pfs)
		case ".tog":
			e, err = tog.New(fe.Name(), pfs)
		case ".ter":
			e, err = ter.New(fe.Name(), pfs)
		//case ".lod":
		//	e, err = lod.New(fe.Name(), pfs)
		case ".prt":
			//e, err = pts.New(fe.Name(), pfs)
			fmt.Println("TODO: prt PTCL support")
			continue
		case ".dds":
			err = os.WriteFile(out+"/"+fe.Name(), fe.Data(), 0644)
			if err != nil {
				return fmt.Errorf("dds.WriteFile %s: %w", fe.Name(), err)
			}
			continue
		case ".txt":
			err = os.WriteFile(out+"/"+fe.Name(), fe.Data(), 0644)
			if err != nil {
				return fmt.Errorf("txt.WriteFile %s: %w", fe.Name(), err)
			}
			continue
		default:
			fmt.Println("TODO:", fe.Name(), "support")
			continue
			//return fmt.Errorf("unknown file type %s", fe.Name())
		}
		if err != nil {
			return fmt.Errorf("%s.New %s: %w", ext, fe.Name(), err)
		}
		err = e.Decode(bytes.NewReader(fe.Data()))
		if err != nil {
			return fmt.Errorf("%s.Decode %s: %w", ext, fe.Name(), err)
		}
		err = e.BlenderExport(out)
		if err != nil {
			return fmt.Errorf("%s.BlenderExport %s: %w", ext, fe.Name(), err)
		}
	}

	fmt.Printf("%s exported in %.1fs\n", filepath.Base(path), time.Since(start).Seconds())
	return nil
}
