package zon

import (
	"bytes"
	"os"
	"testing"

	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
)

func TestGLTF(t *testing.T) {
	category := "bazaar"
	path := "test/eq/" + category + ".eqg"
	archive, err := eqg.New(path)
	if err != nil {
		t.Fatalf("eqg.new: %s", err)
	}
	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("eqg %s", err)
	}
	defer r.Close()
	err = archive.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	zonData, err := archive.File(category + ".zon")
	if err != nil {
		t.Fatalf("file: %s", err)
	}
	e, err := NewEQG(category, archive)
	if err != nil {
		t.Fatalf("newEQG: %s", err)
	}
	err = e.Load(bytes.NewReader(zonData))
	if err != nil {
		t.Fatalf("zon load: %s", err)
	}
	doc, err := gltf.New()
	if err != nil {
		t.Fatalf("gltf.new: %s", err)
	}
	err = e.GLTFExport(doc)
	if err != nil {
		t.Fatalf("gltf: %s", err)
	}
	w, err := os.Create("test/eq/" + category + ".gltf")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = doc.Export(w)
	if err != nil {
		t.Fatalf("export: %s", err)
	}
}
