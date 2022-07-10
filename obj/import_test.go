package obj

import "testing"

func TestImportBox(t *testing.T) {
	req := &ObjRequest{
		ObjPath:    "test/box.obj",
		MtlPath:    "test/box.mtl",
		MattxtPath: "test/box_material.txt",
	}
	err := Import(req)
	if err != nil {
		t.Fatalf("import: %s", err)
	}
	if req.Obj == nil {
		t.Fatalf("empty object")
	}
}

func TestImportRequestNil(t *testing.T) {
	err := Import(nil)
	if err != nil && err.Error() != "request is nil" {
		t.Fatalf("wanted 'request is nil', got %s", err)
	}
}

func TestImportNoPath(t *testing.T) {
	req := &ObjRequest{}
	err := Import(req)
	if err != nil && err.Error() != "importMatTxt: open : no such file or directory" {
		t.Fatalf("wanted 'importMatTxt: open : no such file or directory', got %s", err)
	}
}
