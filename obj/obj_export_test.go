package obj

import (
	"testing"
)

func TestObjExport(t *testing.T) {
	req := &ObjRequest{
		Obj:     &ObjData{},
		ObjPath: "test/tmp.obj",
	}
	err := objExport(req)
	if err != nil {
		t.Fatalf("objExport: %s", err)
	}
}

func TestObjExportRequestNil(t *testing.T) {
	err := objExport(nil)
	if err != nil && err.Error() != "request is nil" {
		t.Fatalf("wanted 'request is nil', got %s", err)
	}
}

func TestObjectExportObjectNil(t *testing.T) {
	req := &ObjRequest{}
	err := objExport(req)
	if err != nil && err.Error() != "request object is nil" {
		t.Fatalf("wanted 'request object is nil', got %s", err)
	}
}
