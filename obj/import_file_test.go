package obj

import (
	"testing"
)

func TestImportObjFile(t *testing.T) {
	req := &ObjRequest{
		Data:    &ObjData{},
		ObjPath: "test/box.obj",
	}
	err := importFile(req)
	if err != nil {
		t.Fatalf("importFile: %s", err)
	}
	//t.Fatalf("%+v", obj)
}
