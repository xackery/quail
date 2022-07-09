package obj

import (
	"testing"
)

func TestImportObjFile(t *testing.T) {
	obj := &ObjData{}
	err := importFile(obj, "test/box.obj")
	if err != nil {
		t.Fatalf("importObj: %s", err)
	}
	//t.Fatalf("%+v", obj)
}
