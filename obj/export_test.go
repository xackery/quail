package obj

import "testing"

func TestExport(t *testing.T) {
	req := &ObjRequest{
		Data:       &ObjData{},
		ObjPath:    "test/tmp.obj",
		MtlPath:    "test/tmp.mtl",
		MattxtPath: "test/tmp_material.txt",
	}
	err := Export(req)
	if err != nil {
		t.Fatalf("Export: %s", err)
	}
	//t.Fatalf("%+v", obj)
}

func TestImportExportBox(t *testing.T) {
	req := &ObjRequest{
		ObjPath:    "test/box.obj",
		MtlPath:    "test/box.mtl",
		MattxtPath: "test/box_material.txt",
	}
	err := Import(req)
	if err != nil {
		t.Fatalf("import: %s", err)
	}
	req.Data.Name = "box"
	req.ObjPath = "test/tmp.obj"
	req.MtlPath = "test/tmp.mtl"
	req.MattxtPath = "test/tmp_material.txt"

	err = Export(req)
	if err != nil {
		t.Fatalf("export: %s", err)
	}
}
