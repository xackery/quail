package mod

import (
	"os"
	"testing"
)

func TestObjImport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	oPath := "../eq/tmp/cube.obj"
	mPath := "../eq/tmp/cube.mtl"
	err = e.ImportObj(oPath, mPath, "")
	if err != nil {
		t.Fatalf("importObj: %s", err)
	}
	w, err := os.Create("../eq/tmp/out.mod")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}
