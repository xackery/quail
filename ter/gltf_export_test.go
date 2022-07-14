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
		{category: "triangle_plane"},
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
