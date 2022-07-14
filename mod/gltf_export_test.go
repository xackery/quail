package mod

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/eqg"
)

func TestGLTFExportSamples(t *testing.T) {
	tests := []struct {
		category string
	}{
		{category: "arthwall"},
		{category: "aro"},
		{category: "she"},
		{category: "voaequip"},
	}
	for _, tt := range tests {

		eqgFile := fmt.Sprintf("test/eq/%s.eqg", tt.category)
		//modFile := fmt.Sprintf("obj_%s.mod", tt.category)

		//outFile := fmt.Sprintf("test/eq/%s_mod.gltf", tt.category)
		//txtFile := fmt.Sprintf("test/eq/%s_mod.txt", tt.category)

		ra, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		defer ra.Close()
		a, err := eqg.New(tt.category)
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		err = a.Load(ra)
		if err != nil {
			t.Fatalf("load eqg: %s", err)
		}

		files := a.Files()
		for _, modEntry := range files {
			if filepath.Ext(modEntry.Name()) != ".mod" {
				continue
			}
			r := bytes.NewReader(modEntry.Data())

			e, err := NewEQG(tt.category, a)
			if err != nil {
				t.Fatalf("new: %s", err)
			}

			err = e.Load(r)
			if err != nil {
				t.Fatalf("load %s: %s", modEntry.Name(), err)
			}

			/*			fw, err := os.Create(txtFile)
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
			outFile := fmt.Sprintf("test/eq/%s_mod_%s.gltf", tt.category, modEntry.Name())
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
}

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
