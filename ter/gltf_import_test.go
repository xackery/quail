package ter

import (
	"testing"
)

func TestGLTFImportExportBoxGLTF(t *testing.T) {
	path := "test/box.eqg"
	inFile := "test/box_out.gltf"

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.GLTFImport(inFile)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}
}
