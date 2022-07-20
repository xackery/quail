package mds

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
	path := "test/eq/lth.eqg"
	inFile := "lth.mds"

	a, err := eqg.New(path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	ra, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer ra.Close()
	err = a.Load(ra)
	if err != nil {
		t.Fatalf("archive load: %s", err)
	}

	dump.New(inFile)
	defer dump.WriteFileClose(path + "_" + inFile + ".png")

	e, err := New(inFile, a)
	if err != nil {
		t.Fatalf("mds new: %s", err)
	}
	data, err := a.File(inFile)
	if err != nil {
		t.Fatalf("file: %s", err)
	}

	err = e.Load(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("load mds: %s", err)
	}
}

func TestLoadSaveLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	filePath := "test/"
	inFile := "test/obj_gears.mod"
	outFile := "test/obj_gears_loadsaveload.mod"

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("newPath: %s", err)
	}

	f, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()

	dump.New(inFile)
	defer dump.WriteFileClose(inFile + ".png")

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
	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("newPath: %s", err)
	}

	f, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()

	dump.New(inFile)
	defer dump.WriteFileClose(inFile + ".png")

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
