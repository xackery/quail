package ter

import (
	"fmt"
	"os"
	"testing"
)

func TestGLTFExportSamples(t *testing.T) {
	tests := []struct {
		category string
	}{
		{category: "heptagon"},
		{category: "hexagon"},
		{category: "pentagon"},
		{category: "octagon"},
		{category: "plane_top_bottom"},
		{category: "plane_top_bottom2"},
		{category: "box"},
		{category: "plane"},
		{category: "triangle"},
		{category: "triangle_material"},
	}
	for _, tt := range tests {
		isGLTFSource := false

		path := "test/" + tt.category + "/_" + tt.category + ".eqg/"
		inFile := "test/" + tt.category + "/_" + tt.category + ".eqg/" + tt.category + ".ter"
		gltfInFile := "test/" + tt.category + "/" + tt.category + ".gltf"
		outFile := "test/" + tt.category + "/tmp.gltf"

		txtFile := "test/" + tt.category + "/" + tt.category + "_ter.txt"

		e, err := New(tt.category, path)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		if isGLTFSource {
			err = e.GLTFImport(gltfInFile)
			if err != nil {
				t.Fatalf("import %s: %s", path, err)
			}
		} else {
			r, err := os.Open(inFile)
			if err != nil {
				t.Fatalf("open %s: %s", path, err)
			}
			defer r.Close()

			err = e.Load(r)
			if err != nil {
				t.Fatalf("load %s: %s", path, err)
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
