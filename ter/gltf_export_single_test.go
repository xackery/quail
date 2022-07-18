package ter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
)

func TestGLTFExportEQGZones(t *testing.T) {
	tests := []struct {
		category string
	}{
		{category: "bazaar"},
		//{category: "steamfontmts"},
		//{category: "broodlands"},
		//{category: "steppes"},
	}
	for _, tt := range tests {
		isDumpEnabled := false

		eqgFile := fmt.Sprintf("test/eq/%s.eqg", tt.category)

		var err error
		var d *dump.Dump
		if isDumpEnabled {
			d, err = dump.New(tt.category)
			if err != nil {
				t.Fatalf("dump.New: %s", err)
			}
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

		for _, fileEntry := range a.Files() {
			if filepath.Ext(fileEntry.Name()) != ".ter" {
				continue
			}

			terBuf := bytes.NewReader(fileEntry.Data())
			err = e.Load(terBuf)
			if err != nil {
				t.Fatalf("load %s: %s", fileEntry.Name(), err)
			}

			w, err := os.Create(fmt.Sprintf("%s.gltf", eqgFile))
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
			if d != nil {
				err = d.Save(fileEntry.Name() + ".png")
				if err != nil {
					t.Fatalf("save png: %s", err)
				}
			}
		}
	}
}

func TestGLTFExportBroodlands(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	zone := "broodlands"
	filePath := fmt.Sprintf("test/eq/_%s.eqg", zone)
	inFile := fmt.Sprintf("test/eq/_%s.eqg/ter_%s.ter", zone, zone)
	outFile := fmt.Sprintf("test/eq/%s.gltf", zone)

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}
	e, err := New("arena", path)
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
		t.Fatalf("import %s: %s", inFile, err)
	}

	/*fw, err := os.Create(fmt.Sprintf("test/%s.txt", zone))
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
		fmt.Fprintf(fw, "%d pos: %+v, normal: %+v, uv: %+v\n", i, o.Position, o.Normal, o.Uv)
	}
	*/
	w, err := os.Create(outFile)
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
}

func TestGLTFExportCityOfBronze(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	zone := "cityofbronze"
	filePath := fmt.Sprintf("test/eq/_%s.eqg", zone)
	inFile := fmt.Sprintf("test/eq/_%s.eqg/ter_%s.ter", zone, zone)
	outFile := fmt.Sprintf("test/eq/%s.gltf", zone)
	isDumpEnabed := false

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}
	if isDumpEnabed {
		d, err := dump.New(path.String())
		if err != nil {
			t.Fatalf("dump.new: %s", err)
		}
		defer d.Save(fmt.Sprintf("%s.png", inFile))
	}

	e, err := New("cityofbronze", path)
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

	/*fw, err := os.Create(fmt.Sprintf("test/%s.txt", zone))
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
		fmt.Fprintf(fw, "%d pos: %+v, normal: %+v, uv: %+v\n", i, o.Position, o.Normal, o.Uv)
	}
	*/
	w, err := os.Create(outFile)
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
}
