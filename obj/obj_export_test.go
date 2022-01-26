package obj

import (
	"os"
	"testing"
)

func TestExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	obj := &ObjData{}
	err := Export(obj, "../eq/tmp/out.obj", "../eq/tmp/out.mtl", "../eq/tmp/out_material.txt")
	if err != nil {
		t.Fatalf("Export: %s", err)
	}
	t.Fatalf("%+v", obj)
}

func TestImportExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	obj, err := Import("../eq/soldungb/cache/soldungb.obj", "../eq/soldungb/cache/soldungb.mtl", "../eq/soldungb/cache/soldungb_material.txt")
	if err != nil {
		t.Fatalf("import: %s", err)
	}
	obj.Name = "soldungb"
	err = Export(obj, "../eq/tmp/out.obj", "../eq/tmp/out.mtl", "../eq/tmp/out_material.txt")
	if err != nil {
		t.Fatalf("Export: %s", err)
	}
	t.Fatalf("%+v", obj)
}
