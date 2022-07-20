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
	archive, err := eqg.NewFile(path)
	if err != nil {
		t.Fatalf("eqg new: %s", err)
	}
	zonData, err := archive.File(category + ".zon")
	if err != nil {
		t.Fatalf("file: %s", err)
	}
	e, err := New(category, archive)
	if err != nil {
		t.Fatalf("newEQG: %s", err)
	}
	err = e.Decode(bytes.NewReader(zonData))
	if err != nil {
		t.Fatalf("zon decode: %s", err)
	}
	doc, err := gltf.New()
	if err != nil {
		t.Fatalf("gltf.new: %s", err)
	}
	err = e.GLTFEncode(doc)
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
