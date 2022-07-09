package obj

import (
	"testing"
)

func TestExportFile(t *testing.T) {
	obj := &ObjData{}
	err := exportFile(obj, "test/tmp.obj")
	if err != nil {
		t.Fatalf("exportObjFile: %s", err)
	}
}
