package obj

import (
	"os"
	"testing"
)

func TestExportObjFile(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	obj := &ObjData{}
	err := exportObjFile(obj, "../eq/tmp/out.obj")
	if err != nil {
		t.Fatalf("exportObjFile: %s", err)
	}
	t.Fatalf("%+v", obj)
}
