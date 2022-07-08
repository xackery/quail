package mod

import (
	"io/fs"
	"os"
	"testing"
)

func TestGLTFExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	err := os.Mkdir("test", fs.ModeDir)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("mkdir test: %s", err)
	}

	e, err := New("obj_gears.mod", "test/")
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open("test/obj_gears.mod")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer r.Close()
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load %s", err)
	}

	w, err := os.Create("test/obj_gears.gltf")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = e.ExportGLTF(w)
	if err != nil {
		t.Fatalf("export: %s", err)
	}
}
