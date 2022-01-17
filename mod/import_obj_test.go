package mod

import (
	"os"
	"testing"
)

func TestObjImport(t *testing.T) {

	e := &MOD{}
	oPath := "test/cube.obj"
	r, err := os.Open(oPath)
	if err != nil {
		t.Fatalf("open %s: %s", oPath, err)
	}
	mPath := "test/cube.mtl"
	mr, err := os.Open(mPath)
	if err != nil {
		t.Fatalf("open %s: %s", mPath, err)
	}
	err = e.ImportObj(r, mr)
	if err != nil {
		t.Fatalf("importObj: %s", err)
	}
	w, err := os.Create("test/out.mod")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}
