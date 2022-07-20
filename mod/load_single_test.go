package mod

import (
	"bytes"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	filePath := "test/eq/_steamfontmts.eqg/"
	inFile := "obj_gears.mod"

	archive, err := eqg.NewFile(filePath)
	if err != nil {
		t.Fatalf("eqg new: %s", err)
	}

	data, err := archive.File(inFile)
	if err != nil {
		t.Fatalf("eqg load: %s", err)
	}

	dump.New(inFile)
	defer dump.WriteFileClose(filePath + "_" + inFile + ".png")

	e, err := New("out", archive)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("load: %s", err)
	}
}

func TestLoadSaveLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	filePath := "test/"
	inFile := "obj_gears.mod"
	outFile := "test/obj_gears_loadsaveload.mod"
	f, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()

	archive, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}

	dump.New(archive.String())
	defer dump.WriteFileClose(filePath + inFile + ".png")

	e, err := NewFile("out", archive, "obj_gears.mod")
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}

	r, err := os.Open(outFile)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
}

func TestLoadSaveGLTF(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	filePath := "test/"
	inFile := "test/obj_gears.mod"
	outFile := "test/obj_gears_loadsavegtlf.gltf"

	f, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()

	dump.New(inFile)
	defer dump.WriteFileClose(inFile)

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}
	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create gltf: %s", err)
	}
	defer w.Close()

	doc, err := gltf.New()
	if err != nil {
		t.Fatalf("gltf.New: %s", err)
	}
	err = e.GLTFExport(doc)
	if err != nil {
		t.Fatalf("gltf: %s", err)
	}
	err = doc.Export(w)
	if err != nil {
		t.Fatalf("export: %s", err)
	}
}
