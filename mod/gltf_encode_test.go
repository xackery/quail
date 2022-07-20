package mod

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
	}
	for _, tt := range tests {

		filePath := "test/"
		gltfInFile := fmt.Sprintf("test/%s.gltf", tt.category)
		gltfOutFile := fmt.Sprintf("test/%s_out.gltf", tt.category)

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
			t.Fatalf("gltf: %s", err)
		}

		err = doc.Export(w)
		if err != nil {
			t.Fatalf("export: %s", err)
		}
	}
}
