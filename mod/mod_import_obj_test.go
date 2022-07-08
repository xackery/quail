package mod

import (
	"os"
	"testing"
)

func TestObjImport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/"
	inFileObj := "test/cube.obj"
	inFileMat := "test/cube.mtl"
	outFile := "test/cube_objimport.mod"

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.ImportObj(inFileObj, inFileMat, "")
	if err != nil {
		t.Fatalf("importObj: %s", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}
