package mds

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xackery/quail/model/metadata/lay"
	"github.com/xackery/quail/model/plugin/gltf"
	"github.com/xackery/quail/pfs/eqg"
)

func TestGLTFFlushEQPath(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	files, err := os.ReadDir("test/eq/")
	if err != nil {
		t.Fatalf("readdir: %s", err)
	}
	for _, fe := range files {
		if !strings.HasSuffix(fe.Name(), ".gltf") {
			continue
		}
		err = os.Remove(fmt.Sprintf("test/eq/%s", fe.Name()))
		if err != nil {
			t.Fatalf("remove: %s", err)
		}
	}
}

func TestGLTFEncodeSamplesSingleTest(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	tests := []struct {
		category string
	}{
		/*{category: "bxi"},
		{category: "bnf"},
		{category: "bnm"},
		{category: "bnr"},
		{category: "bnx"},*/
		//{category: "bny"},
		//{category: "dv6"},
		//{category: "lth"}, //Morell thule
		{category: "djm"}, //Djinn Male
		//{category: "zmm"}, //Zombie Male
		//{category: "wrm"}, //worm
		{category: "dgj"}, //jade dragon
		//{category: "prt"},
	}
	for _, tt := range tests {

		fmt.Println("loading", tt.category)
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
		err = a.Decode(ra)
		if err != nil {
			t.Fatalf("decode eqg: %s", err)
		}

		files := a.Files()
		for _, mdsEntry := range files {
			if filepath.Ext(mdsEntry.Name()) != ".mds" {
				continue
			}
			r := bytes.NewReader(mdsEntry.Data())

			e, err := New(strings.TrimSuffix(mdsEntry.Name(), ".mds"), a)
			if err != nil {
				t.Fatalf("mds new: %s", err)
			}

			err = e.Decode(r)
			if err != nil {
				t.Fatalf("decode %s: %s", mdsEntry.Name(), err)
			}

			layName := fmt.Sprintf("%s.lay", strings.TrimSuffix(mdsEntry.Name(), ".mds"))
			layEntry, err := a.File(layName)
			if err != nil && !strings.Contains(err.Error(), "does not exist") {
				t.Fatalf("file: %s", err)
			}

			if len(layEntry) > 0 {
				l, err := lay.New(layName, a)
				if err != nil {
					t.Fatalf("lay.NewEQG: %s", err)
				}
				err = l.Decode(bytes.NewReader(layEntry))
				if err != nil {
					t.Fatalf("decode lay: %s", err)
				}
				err = e.SetLayers(l.Layers())
				if err != nil {
					t.Fatalf("setLayers: %s", err)
				}
			}

			outFile := fmt.Sprintf("test/eq/%s_eqg_%s.gltf", tt.category, mdsEntry.Name())
			w, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("create %s", err)
			}
			defer w.Close()
			doc, err := gltf.New()
			if err != nil {
				t.Fatalf("gltf.New: %s", err)
			}
			err = e.GLTFEncode(doc)
			if err != nil {
				t.Fatalf("gltf: %s", err)
			}
			err = doc.Export(w)
			if err != nil {
				t.Fatalf("export: %s", err)
			}

		}
	}
}

func TestGLTFEncodeSingleModel(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	tests := []struct {
		category string
		model    string
	}{
		{category: "steamfontmts", model: "obj_gears.mod"},
		{category: "steamfontmts", model: "obj_oilbarrel.mod"},
		{category: "voaequip", model: "obj_oilbarrel.mod"},
		//{category: "arthwall"},
		//{category: "aro"},
		//{category: "she"},
		//{category: "voaequip"},
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
		err = a.Decode(ra)
		if err != nil {
			t.Fatalf("decode eqg: %s", err)
		}

		data, err := a.File(tt.model)
		if err != nil {
			t.Fatalf("file: %s", err)
		}

		e, err := New(tt.model, a)
		if err != nil {
			t.Fatalf("mds new: %s", err)
		}

		r := bytes.NewReader(data)
		err = e.Decode(r)
		if err != nil {
			t.Fatalf("decode %s: %s", tt.model, err)
		}

		outFile := fmt.Sprintf("test/eq/%s_mod_%s.gltf", tt.category, tt.model)
		w, err := os.Create(outFile)
		if err != nil {
			t.Fatalf("create %s", err)
		}
		defer w.Close()
		doc, err := gltf.New()
		if err != nil {
			t.Fatalf("gltf.New: %s", err)
		}
		err = e.GLTFEncode(doc)
		if err != nil {
			t.Fatalf("gltf: %s", err)
		}

		err = doc.Export(w)
		if err != nil {
			t.Fatalf("export: %s", err)
		}
	}
}

func TestGLTFEncodeSamplesDirectory(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	dirs, err := os.ReadDir("test/eq/")
	if err != nil {
		t.Fatalf("read: %s", err)
	}
	for _, fe := range dirs {
		if filepath.Ext(fe.Name()) != ".eqg" {
			continue
		}
		if fe.IsDir() {
			continue
		}
		fmt.Println("loading", fe.Name())
		eqgFile := fmt.Sprintf("test/eq/%s", fe.Name())

		ra, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		defer ra.Close()
		a, err := eqg.New(fe.Name())
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		err = a.Decode(ra)
		if err != nil {
			t.Fatalf("decode eqg: %s", err)
		}

		files := a.Files()
		for _, mdsEntry := range files {
			if filepath.Ext(mdsEntry.Name()) != ".mds" {
				continue
			}
			if fe.IsDir() {
				continue
			}
			r := bytes.NewReader(mdsEntry.Data())

			e, err := New(strings.TrimSuffix(mdsEntry.Name(), ".mds"), a)
			if err != nil {
				t.Fatalf("mds new: %s", err)
			}

			err = e.Decode(r)
			if err != nil {
				t.Fatalf("decode %s: %s", mdsEntry.Name(), err)
			}

			layName := fmt.Sprintf("%s.lay", strings.TrimSuffix(mdsEntry.Name(), ".mds"))
			layEntry, err := a.File(layName)
			if err != nil && !strings.Contains(err.Error(), "does not exist") {
				t.Fatalf("file: %s", err)
			}

			if len(layEntry) > 0 {
				l, err := lay.New(layName, a)
				if err != nil {
					t.Fatalf("lay.NewEQG: %s", err)
				}
				err = l.Decode(bytes.NewReader(layEntry))
				if err != nil {
					t.Fatalf("decode lay: %s", err)
				}
				err = e.SetLayers(l.Layers())
				if err != nil {
					t.Fatalf("setLayers: %s", err)
				}
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
			outFile := fmt.Sprintf("test/eq/%s_eqg_%s.gltf", fe.Name(), mdsEntry.Name())
			w, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("create %s", err)
			}
			defer w.Close()
			doc, err := gltf.New()
			if err != nil {
				t.Fatalf("gltf.New: %s", err)
			}
			err = e.GLTFEncode(doc)
			if err != nil {
				t.Fatalf("gltf: %s", err)
			}

			err = doc.Export(w)
			if err != nil {
				t.Fatalf("export: %s", err)
			}
		}
	}
}

func TestGLTFEncodeMaterialDump(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	dirs, err := os.ReadDir("test/eq/")
	if err != nil {
		t.Fatalf("read: %s", err)
	}
	for _, fe := range dirs {
		if filepath.Ext(fe.Name()) != ".eqg" {
			continue
		}
		if fe.IsDir() {
			continue
		}
		eqgFile := fmt.Sprintf("test/eq/%s", fe.Name())

		ra, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		defer ra.Close()
		archive, err := eqg.New(fe.Name())
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		err = archive.Decode(ra)
		if err != nil {
			t.Fatalf("decode eqg: %s", err)
		}

		files := archive.Files()
		for _, mdsEntry := range files {
			if filepath.Ext(mdsEntry.Name()) != ".mds" {
				continue
			}
			r := bytes.NewReader(mdsEntry.Data())

			e, err := New(strings.TrimSuffix(mdsEntry.Name(), ".mds"), archive)
			if err != nil {
				t.Fatalf("mds new: %s", err)
			}

			err = e.Decode(r)
			if err != nil {
				t.Fatalf("decode %s: %s", mdsEntry.Name(), err)
			}

			layName := fmt.Sprintf("%s.lay", strings.TrimSuffix(mdsEntry.Name(), ".mds"))
			layEntry, err := archive.File(layName)
			if err != nil && !strings.Contains(err.Error(), "does not exist") {
				t.Fatalf("file: %s", err)
			}

			if len(layEntry) > 0 {
				l, err := lay.New(layName, archive)
				if err != nil {
					t.Fatalf("lay.NewEQG: %s", err)
				}
				err = l.Decode(bytes.NewReader(layEntry))
				if err != nil {
					t.Fatalf("decode lay: %s", err)
				}
				err = e.SetLayers(l.Layers())
				if err != nil {
					t.Fatalf("setLayers: %s", err)
				}
			}
			for _, mat := range e.materials {
				fmt.Println(mdsEntry.Name(), mat.Name)
			}
		}
	}
}
