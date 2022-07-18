package ter

import (
	"os"
	"testing"

	"github.com/xackery/quail/common"
)

func TestObjExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	filePath := "test/"
	inFile := "test/soldungb.ter"

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	err = e.ObjExport("test/objexport.obj", "test/objexport.mtl", "test/objexport.txt")
	if err != nil {
		t.Fatalf("export: %s", err)
	}

}
