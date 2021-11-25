package s3d

import (
	"os"
	"testing"
)

func TestWldImportGltf(t *testing.T) {
	e := &Wld{}
	f, err := os.Open("test/box.gltf")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	err = e.ImportGltf(f)
	if err != nil {
		t.Fatalf("importgltf: %v", err)
	}
}
