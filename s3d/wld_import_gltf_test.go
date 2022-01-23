package s3d

import (
	"os"
	"testing"
)

func TestWldImportGltf(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &Wld{}
	f, err := os.Open("../eq/tmp/box.gltf")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	err = e.ImportGltf(f)
	if err != nil {
		t.Fatalf("importgltf: %v", err)
	}
}
