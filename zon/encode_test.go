package zon

import (
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
	"github.com/xackery/quail/ter"
)

func TestEncodeEQG(t *testing.T) {
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

	e, err := ter.New(category, path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	gdoc, err := gltf.Open("test/monkey.gltf")
	if err != nil {
		t.Fatalf("gltf open: %s", err)
	}
	err = e.GLTFDecode(gdoc)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	err = e.ArchiveExport(archive)
	if err != nil {
		t.Fatalf("archive export: %s", err)
	}

	zone, err := New(category, archive)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	zone.models = append(zone.models, &model{
		name:     category + ".ter",
		baseName: category,
	})

	err = zone.ArchiveExport(archive)
	if err != nil {
		t.Fatalf("zon archive export: %s", err)
	}

	w, err := os.Create("test/" + category + ".eqg")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = archive.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}

}
