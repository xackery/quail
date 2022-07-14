package ter

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/eqg"
)

func TestGLTFExportSamples(t *testing.T) {
	tests := []struct {
		category string
	}{
		{category: "box"},
		{category: "heptagon"},
		{category: "hexagon"},
		{category: "octagon"},
		{category: "pentagon"},
		{category: "plane_top_bottom"},
		{category: "plane_top_bottom2"},
		{category: "plane"},
		{category: "triangle_material"},
		{category: "triangle"},
	}
	for _, tt := range tests {
		isGLTFSource := false

		eqgFile := fmt.Sprintf("test/%s.eqg", tt.category)
		terFile := fmt.Sprintf("%s.ter", tt.category)
		gltfFile := fmt.Sprintf("test/%s.gltf", tt.category)
		gltfOutFile := fmt.Sprintf("test/%s_out.gltf", tt.category)
		txtFile := fmt.Sprintf("test/%s_ter.txt", tt.category)
		if isGLTFSource {
			txtFile = fmt.Sprintf("test/%s_gltf.txt", tt.category)
		}

		a, err := eqg.New(tt.category)
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		r, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		err = a.Load(r)
		if err != nil {
			t.Fatalf("load: %s", err)
		}

		e, err := NewEQG(tt.category, a)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		if isGLTFSource {
			err = e.GLTFImport(gltfFile)
			if err != nil {
				t.Fatalf("import %s: %s", gltfFile, err)
			}
		} else {
			data, err := a.File(terFile)
			if err != nil {
				t.Fatalf("File: %s", err)
			}
			r := bytes.NewReader(data)
			err = e.Load(r)
			if err != nil {
				t.Fatalf("load %s: %s", terFile, err)
			}
		}

		fw, err := os.Create(txtFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		defer fw.Close()
		fmt.Fprintf(fw, "faces:\n")
		for i, o := range e.faces {
			fmt.Fprintf(fw, "%d %+v\n", i, o)
		}

		fmt.Fprintf(fw, "vertices:\n")
		for i, o := range e.vertices {
			fmt.Fprintf(fw, "%d pos: %0.0f %0.0f %0.0f, normal: %+v, uv: %+v\n", i, o.Position.X, o.Position.Y, o.Position.Z, o.Normal, o.Uv)
		}

		w, err := os.Create(gltfOutFile)
		if err != nil {
			t.Fatalf("create %s", err)
		}
		defer w.Close()
		err = e.GLTFExport(w)
		if err != nil {
			t.Fatalf("save: %s", err)
		}
	}
}

func TestGLTFExportBoxFromGLTF(t *testing.T) {
	path := "test/box/"
	inFile := "test/box/_box.eqg/box.gltf"
	outFile := "test/box/tmp.gltf"

	e, err := New("box", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.GLTFImport(inFile)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	fw, err := os.Create("test/box_gltf.txt")
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
		fmt.Fprintf(fw, "%d pos: %0.0f %0.0f %0.0f, normal: %+v, uv: %+v\n", i, o.Position.X, o.Position.Y, o.Position.Z, o.Normal, o.Uv)
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

func TestGLTFExportTriangleFromGLTF(t *testing.T) {
	path := "test/triangle/"
	inFile := "test/triangle/_triangle.eqg/triangle.gltf"
	outFile := "test/triangle/tmp.gltf"

	e, err := New("triangle", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.GLTFImport(inFile)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	fw, err := os.Create("test/triangle_gltf.txt")
	if err != nil {
		t.Fatalf("triangle.txt: %s", err)
	}
	defer fw.Close()
	fmt.Fprintf(fw, "faces:\n")
	for i, o := range e.faces {
		fmt.Fprintf(fw, "%d %+v\n", i, o)
	}

	fmt.Fprintf(fw, "vertices:\n")
	for i, o := range e.vertices {
		fmt.Fprintf(fw, "%d pos: %0.0f %0.0f %0.0f, normal: %+v, uv: %+v\n", i, o.Position.X, o.Position.Y, o.Position.Z, o.Normal, o.Uv)
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
