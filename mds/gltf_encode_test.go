package mds

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/gltf"
)

func TestGLTFEncodeGLTF(t *testing.T) {
	tests := []struct {
		category string
	}{
		{category: "animation"},
		//{category: "simple"},
		//{category: "man"},
	}
	for _, tt := range tests {

		path := "test/"
		p, err := common.NewPath(path)
		if err != nil {
			t.Fatalf("newPath: %s", err)
		}
		gltfInFile := fmt.Sprintf("test/%s.gltf", tt.category)
		gltfOutFile := fmt.Sprintf("test/%s_out.gltf", tt.category)
		e, err := New(tt.category, p)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		gdoc, err := gltf.Open(gltfInFile)
		if err != nil {
			t.Fatalf("gltf open: %s", err)
		}
		err = e.GLTFDecode(gdoc)
		if err != nil {
			t.Fatalf("gltfimport '%s': %s", tt.category, err)
		}

		w, err := os.Create(gltfOutFile)
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
			t.Fatalf("gltfExport: %s", err)
		}

		err = doc.Export(w)
		if err != nil {
			t.Fatalf("export: %s", err)
		}
	}
}
