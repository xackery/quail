package prt

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
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

func TestGLTFExportSamplesSingleTest(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	tests := []struct {
		category string
	}{
		{category: "djm"},
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
		archive, err := eqg.New(tt.category)
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		err = archive.Load(ra)
		if err != nil {
			t.Fatalf("load eqg: %s", err)
		}

		files := archive.Files()
		for _, modEntry := range files {
			if filepath.Ext(modEntry.Name()) != ".prt" {
				continue
			}
			r := bytes.NewReader(modEntry.Data())

			e, err := New(modEntry.Name(), archive)
			if err != nil {
				t.Fatalf("new: %s", err)
			}

			err = e.Load(r)
			if err != nil {
				t.Fatalf("load %s: %s", modEntry.Name(), err)
			}

			outFile := fmt.Sprintf("test/eq/%s_eqg_%s.gltf", tt.category, modEntry.Name())
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
	}
}
