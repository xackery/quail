package mod

import (
	"os"
	"testing"

	"github.com/xackery/quail/pfs/archive"
)

func TestObjImport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	filePath := "test/"
	inFileObj := "test/box/cache/box.obj"
	inFileMat := "test/box/cache/box.mtl"
	outFile := "test/box.mod"

	path, err := archive.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}
	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.ObjImport(inFileObj, inFileMat, "")
	if err != nil {
		t.Fatalf("importObj: %s", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}
}
