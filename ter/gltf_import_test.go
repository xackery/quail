package ter

import (
	"testing"

	"github.com/xackery/quail/eqg"
	"github.com/xackery/quail/gltf"
)

func TestGLTFImportExportBoxGLTF(t *testing.T) {
	filePath := "test/box.eqg"
	inFile := "test/box_out.gltf"

	archive, err := eqg.New(filePath)
	if err != nil {
		t.Fatalf("eqg new : %s", err)
	}
	e, err := New("out", archive)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	gdoc, err := gltf.Open(inFile)
	if err != nil {
		t.Fatalf("gltf open: %s", err)
	}
	err = e.GLTFImport(gdoc)
	if err != nil {
		t.Fatalf("import %s: %s", filePath, err)
	}
}
