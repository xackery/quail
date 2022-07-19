package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/eqg"
	qexport "github.com/xackery/quail/export"
	"github.com/xackery/quail/gltf"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
	"github.com/xackery/quail/zon"
)

func viewLoad(buf *bytes.Buffer, path string, file string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	ext := strings.ToLower(filepath.Ext(path))

	type viewFunc struct {
		name   string
		invoke func(buf *bytes.Buffer, path string, file string, ext string) error
	}
	views := []*viewFunc{
		{name: "eqg", invoke: viewLoadEQG},
		{name: "gltf", invoke: viewLoadGLTF},
	}

	for _, v := range views {
		err = v.invoke(buf, path, file, ext)
		if err != nil {
			return fmt.Errorf("view load %s: %w", v.name, err)
		}
		if buf.Len() > 0 {
			return nil
		}
	}
	return fmt.Errorf("%s is not viewable", path)
}

func viewLoadGLTF(buf *bytes.Buffer, path string, file string, ext string) error {
	if ext != ".gltf" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = buf.Write(data)
	if err != nil {
		return fmt.Errorf("buf.Write: %w", err)
	}
	return nil
}

func viewLoadEQG(buf *bytes.Buffer, path string, file string, ext string) error {
	if ext != ".eqg" {
		return nil
	}
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	archive, err := eqg.New(path)
	if err != nil {
		return fmt.Errorf("eqg.New: %w", err)
	}
	err = archive.Load(r)
	if err != nil {
		return fmt.Errorf("eqg load: %w", err)
	}

	doc, err := gltf.New()
	if err != nil {
		return fmt.Errorf("gltf new: %w", err)
	}

	if file != "" {
		data, err := archive.File(file)
		if err != nil {
			return fmt.Errorf("eqg file: %w", err)
		}
		_, err = buf.Write(data)
		if err != nil {
			return fmt.Errorf("buf.Write: %w", err)
		}
	} else {
		e, err := qexport.New(filepath.Base(path), archive)
		if err != nil {
			return fmt.Errorf("export new: %w", err)
		}

		err = e.LoadArchive()
		if err != nil {
			return fmt.Errorf("load archive: %w", err)
		}

		err = e.GLTFExport(doc)
		if err != nil {
			return fmt.Errorf("gltfexport: %w", err)
		}

		err = doc.Export(buf)
		if err != nil {
			return fmt.Errorf("gltf export: %w", err)
		}
	}
	return nil
}

func viewLoadMDS(buf *bytes.Buffer, path string, file string, ext string) error {
	if ext != ".mds" {
		return nil
	}
	archive, err := common.NewPath(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("path new: %w", err)
	}

	e, err := mds.New(filepath.Base(path), archive)
	if err != nil {
		return fmt.Errorf("mds new: %w", err)
	}

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		return fmt.Errorf("mds load: %w", err)
	}

	doc, err := gltf.New()
	if err != nil {
		return fmt.Errorf("gltf new: %w", err)
	}

	err = e.GLTFExport(doc)
	if err != nil {
		return fmt.Errorf("gltfexport: %w", err)
	}

	err = doc.Export(buf)
	if err != nil {
		return fmt.Errorf("gltf export: %w", err)
	}

	return nil
}

func viewLoadMOD(buf *bytes.Buffer, path string, file string, ext string) error {
	if ext != ".mod" {
		return nil
	}
	archive, err := common.NewPath(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("path new: %w", err)
	}

	e, err := mod.New(filepath.Base(path), archive)
	if err != nil {
		return fmt.Errorf("mod new: %w", err)
	}

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		return fmt.Errorf("mod load: %w", err)
	}

	doc, err := gltf.New()
	if err != nil {
		return fmt.Errorf("gltf new: %w", err)
	}

	err = e.GLTFExport(doc)
	if err != nil {
		return fmt.Errorf("gltfexport: %w", err)
	}

	err = doc.Export(buf)
	if err != nil {
		return fmt.Errorf("gltf export: %w", err)
	}

	return nil
}

func viewLoadTER(buf *bytes.Buffer, path string, file string, ext string) error {
	if ext != ".ter" {
		return nil
	}
	archive, err := common.NewPath(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("path new: %w", err)
	}

	e, err := ter.New(filepath.Base(path), archive)
	if err != nil {
		return fmt.Errorf("ter new: %w", err)
	}

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		return fmt.Errorf("ter load: %w", err)
	}

	doc, err := gltf.New()
	if err != nil {
		return fmt.Errorf("gltf new: %w", err)
	}

	err = e.GLTFExport(doc)
	if err != nil {
		return fmt.Errorf("gltfexport: %w", err)
	}

	err = doc.Export(buf)
	if err != nil {
		return fmt.Errorf("gltf export: %w", err)
	}

	return nil
}

func viewLoadZON(buf *bytes.Buffer, path string, file string, ext string) error {
	if ext != ".zon" {
		return nil
	}
	archive, err := common.NewPath(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("path new: %w", err)
	}

	e, err := zon.New(filepath.Base(path), archive)
	if err != nil {
		return fmt.Errorf("zon new: %w", err)
	}

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		return fmt.Errorf("zon load: %w", err)
	}

	doc, err := gltf.New()
	if err != nil {
		return fmt.Errorf("gltf new: %w", err)
	}

	err = e.GLTFExport(doc)
	if err != nil {
		return fmt.Errorf("gltfexport: %w", err)
	}

	err = doc.Export(buf)
	if err != nil {
		return fmt.Errorf("gltf export: %w", err)
	}

	return nil
}
