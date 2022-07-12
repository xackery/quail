package mod

import (
	"io/fs"
	"io/ioutil"
	"os"
	"testing"

	"github.com/xackery/quail/common"
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

	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("export: %s", err)
	}
}

func TestTriangle(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/"
	inFile := "test/triangle.gltf"
	outFile := "test/triangle_out.gltf"

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.GLTFImport(inFile)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	e.materials = append(e.materials, &common.Material{Name: "metal_rustyb.dds", Properties: common.Properties{{Name: "e_texturediffuse0", Value: "metal_rustyb.dds", Category: 2}}})
	data, err := ioutil.ReadFile("test/metal_rustyb.dds")
	if err != nil {
		t.Fatalf("%s", err)
	}
	fe, err := common.NewFileEntry("metal_rustyb.dds", data)
	if err != nil {
		t.Fatalf("NewFileEntry: %s", err)
	}
	e.files = append(e.files, fe)
	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.GLTFExport(w)
	//err = e.Save(w)
	if err != nil {
		t.Fatalf("gltfExport: %s", err)
	}
}
