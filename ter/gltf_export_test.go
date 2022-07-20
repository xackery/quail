package ter

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
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
		isGLTFSource := true
		isDumpEnabled := false

		eqgFile := fmt.Sprintf("test/%s.eqg", tt.category)
		terFile := fmt.Sprintf("%s.ter", tt.category)
		gltfFile := fmt.Sprintf("test/%s.gltf", tt.category)
		gltfOutFile := fmt.Sprintf("test/%s_out.gltf", tt.category)

		var err error

		if isDumpEnabled {
			dump.New(tt.category)
			dump.WriteFileClose(fmt.Sprintf("test/%s_eqg_%s", tt.category, tt.category))
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

		e, err := New(tt.category, a)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		if isGLTFSource {
			gdoc, err := gltf.Open(gltfFile)
			if err != nil {
				t.Fatalf("gltf open: %s", err)
			}
			err = e.GLTFImport(gdoc)
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

		/*fw, err := os.Create(txtFile)
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
		*/
		w, err := os.Create(gltfOutFile)
		if err != nil {
			t.Fatalf("create %s", err)
		}
		defer w.Close()

		doc, err := gltf.New()
		if err != nil {
			t.Fatalf("gltf.New: %s", err)
		}
		err = e.GLTFExport(doc)
		if err != nil {
			t.Fatalf("gltf: %s", err)
		}

		err = doc.Export(w)
		if err != nil {
			t.Fatalf("export: %s", err)
		}
		dump.WriteFileClose(fmt.Sprintf("test/%s_eqg_%s", tt.category, tt.category))
	}
}

func TestGLTFExportSamplesSanityCheck(t *testing.T) {
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

		gltfInFile := fmt.Sprintf("test/%s.gltf", tt.category)
		gltfOutFile := fmt.Sprintf("test/%s_out.gltf", tt.category)

		filePath := "test/"
		path, err := common.NewPath(filePath)
		if err != nil {
			t.Fatalf("path: %s", err)
		}
		e, err := New(tt.category, path)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		gdoc, err := gltf.Open(gltfInFile)
		if err != nil {
			t.Fatalf("gltf open: %s", err)
		}
		err = e.GLTFImport(gdoc)
		if err != nil {
			t.Fatalf("import %s: %s", gltfInFile, err)
		}
		/*
			fw, err := os.Create(fmt.Sprintf("%s.txt", gltfInFile))
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
		*/

		w, err := os.Create(gltfOutFile)
		if err != nil {
			t.Fatalf("create %s", err)
		}
		defer w.Close()

		doc, err := gltf.New()
		if err != nil {
			t.Fatalf("gltf.New: %s", err)
		}
		err = e.GLTFExport(doc)
		if err != nil {
			t.Fatalf("gltf: %s", err)
		}

		err = doc.Export(w)
		if err != nil {
			t.Fatalf("export: %s", err)
		}

		e, err = New(tt.category, path)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		gdoc, err = gltf.Open(gltfOutFile)
		if err != nil {
			t.Fatalf("gltf open: %s", err)
		}
		err = e.GLTFImport(gdoc)
		if err != nil {
			t.Fatalf("import %s: %s", gltfInFile, err)
		}
		/*
			fw2, err := os.Create(fmt.Sprintf("%s.txt", gltfOutFile))
			if err != nil {
				t.Fatalf("%s", err)
			}
			defer fw2.Close()
			fmt.Fprintf(fw2, "faces:\n")
			for i, o := range e.faces {
				fmt.Fprintf(fw2, "%d %+v\n", i, o)
			}

			fmt.Fprintf(fw2, "vertices:\n")
			for i, o := range e.vertices {
				fmt.Fprintf(fw2, "%d pos: %0.0f %0.0f %0.0f, normal: %+v, uv: %+v\n", i, o.Position.X, o.Position.Y, o.Position.Z, o.Normal, o.Uv)
			}
		*/
	}
}
