package ter

import (
	"os"
	"testing"
)

func TestObjExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &TER{}
	path := "../eq/_soldungb.eqg/soldungb.ter"
	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r, "soldungb")
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	err = e.ExportObj("../eq/tmp/out.obj", "../eq/tmp/out.mtl", "../eq/tmp/out_material.txt")
	if err != nil {
		t.Fatalf("export: %s", err)
	}

}
