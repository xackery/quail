package mod

import (
	"os"
	"testing"
)

func TestObjImport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &MOD{}
	oPath := "../eq/tmp/cube.obj"
	r, err := os.Open(oPath)
	if err != nil {
		t.Fatalf("open %s: %s", oPath, err)
	}
	mPath := "../eq/tmp/cube.mtl"
	mr, err := os.Open(mPath)
	if err != nil {
		t.Fatalf("open %s: %s", mPath, err)
	}
	err = e.ImportObj(r, mr)
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
