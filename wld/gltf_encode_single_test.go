package wld

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

func TestGLTFEncodeES3Dones(t *testing.T) {
	tests := []struct {
		category string
	}{
		{category: "crushbone"},
		//{category: "steamfontmts"},
		//{category: "broodlands"},
		//{category: "steppes"},
	}
	for _, tt := range tests {
		isDumpEnabled := false

		eqgFile := fmt.Sprintf("test/eq/%s.s3d", tt.category)

		var err error

		a, err := eqg.New(tt.category)
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		r, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		err = a.Decode(r)
		if err != nil {
			t.Fatalf("decode: %s", err)
		}

		e, err := New(tt.category, a)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		for _, fileEntry := range a.Files() {
			if filepath.Ext(fileEntry.Name()) != ".wld" {
				continue
			}
			if isDumpEnabled {
				dump.New(fmt.Sprintf("%s_%s", eqgFile, fileEntry.Name()))
				dump.WriteFileClose(fmt.Sprintf("test/%s_eqg_%s", tt.category, fileEntry.Name()))
			}
			terBuf := bytes.NewReader(fileEntry.Data())
			err = e.Decode(terBuf)
			if err != nil {
				t.Fatalf("decode %s: %s", fileEntry.Name(), err)
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
			err = e.GLTFEncode(doc)
			if err != nil {
				t.Fatalf("gltf: %s", err)
			}

			err = doc.Export(w)
			if err != nil {
				t.Fatalf("export: %s", err)
			}
			dump.WriteFileClose(fmt.Sprintf("test/%s_eqg_%s", tt.category, fileEntry.Name()))
		}
	}
}

func TestGLTFEncodeBroodlands(t *testing.T) {
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

	err = e.Decode(r)
	if err != nil {
		t.Fatalf("import %s: %s", inFile, err)
	}

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

func TestGLTFEncodeCityOfBronze(t *testing.T) {
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
		dump.New(path.String())
		defer dump.WriteFileClose(inFile)
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

	err = e.Decode(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

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
