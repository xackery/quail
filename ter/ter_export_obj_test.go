package ter

import (
	"os"
	"testing"
)

func TestObjExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	path := "../eq/_soldungb.eqg/soldungb.ter"
	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	err = e.ExportObj("../eq/tmp/out.obj", "../eq/tmp/out.mtl", "../eq/tmp/out_material.txt")
	if err != nil {
		t.Fatalf("export: %s", err)
	}

}
