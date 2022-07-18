package ter

import (
	"testing"

	"github.com/xackery/quail/common"
)

func TestGLTFImportExportBoxGLTF(t *testing.T) {
	filePath := "test/box.eqg"
	inFile := "test/box_out.gltf"

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}
	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.GLTFImport(inFile)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}
}
