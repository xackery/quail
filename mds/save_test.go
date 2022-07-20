package mds

import (
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
)

func TestSaveEQG(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	var err error

	category := "arena"

	path, err := common.NewPath("test/")
	if err != nil {
		t.Fatalf("newPath: %s", err)
	}

	archive, err := eqg.New(category)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	e, err := New(category, path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	gdoc, err := gltf.Open("test/monkey.gltf")
	if err != nil {
		t.Fatalf("gltf open: %s", err)
	}
	err = e.GLTFImport(gdoc)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	err = e.ArchiveExport(archive)
	if err != nil {
		t.Fatalf("archive export: %s", err)
	}

	w, err := os.Create("test/" + category + ".eqg")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = archive.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}

}
