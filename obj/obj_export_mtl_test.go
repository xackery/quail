package obj

import (
	"os"
	"testing"
)

func TestExportMtl(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	obj := &ObjData{}
	err := exportMtl(obj, "../eq/tmp/out.mtl")
	if err != nil {
		t.Fatalf("exportMtl: %s", err)
	}
}
