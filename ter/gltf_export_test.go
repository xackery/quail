package ter

import (
	"fmt"
	"os"
	"testing"
)

func TestGLTFExportTriangleMaterial(t *testing.T) {
	path := "test/triangle_material/_triangle_material.eqg/"
	inFile := "test/triangle_material/_triangle_material.eqg/triangle_material.ter"
	//inFile := "test/ecommons.ter"
	outFile := "test/triangle_material/tmp.gltf"

	e, err := New("triangle_material", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create %s", err)
	}
	defer w.Close()
	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestGLTFExportTriangle(t *testing.T) {
	path := "test/triangle/_triangle.eqg"
	//inFile := "test/_triangle.eqg/triangle.ter"
	inFile := "test/triangle/_triangle.eqg/triangle.ter"
	outFile := "test/triangle/tmp.gltf"

	e, err := New("triangle", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create %s", err)
	}
	defer w.Close()
	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestGLTFExportBox(t *testing.T) {
	path := "test/box/"
	inFile := "test/box/_box.eqg/box.ter"
	outFile := "test/box/tmp.gltf"

	e, err := New("box", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()

	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	fw, err := os.Create("test/box.txt")
	if err != nil {
		t.Fatalf("box.txt: %s", err)
	}
	defer fw.Close()
	fmt.Fprintf(fw, "faces:\n")
	for i, o := range e.faces {
		fmt.Fprintf(fw, "%d %+v\n", i, o)
	}

	fmt.Fprintf(fw, "vertices:\n")
	for i, o := range e.vertices {
		fmt.Fprintf(fw, "%d pos: %+v, normal: %+v, uv: %+v\n", i, o.Position, o.Normal, o.Uv)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create %s", err)
	}
	defer w.Close()
	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}
