package mds

import (
	"fmt"
	"os"
	"testing"
)

func TestGLTFExportGLTF(t *testing.T) {
	tests := []struct {
		category string
		model    string
	}{
		//{category: "animations"},
		{category: "animation"},
	}
	for _, tt := range tests {

		path := "test/"
		gltfInFile := fmt.Sprintf("test/%s.gltf", tt.category)
		gltfOutFile := fmt.Sprintf("test/%s_out.gltf", tt.category)
		e, err := New(tt.model, path)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		err = e.GLTFImport(gltfInFile)
		if err != nil {
			t.Fatalf("gltfiimport %s: %s", tt.model, err)
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
