package obj

import "testing"

func TestExport(t *testing.T) {
	obj := &ObjData{}
	err := Export(obj, "test/tmp.obj", "test/tmp.mtl", "test/tmp_material.txt")
	if err != nil {
		t.Fatalf("Export: %s", err)
	}
	//t.Fatalf("%+v", obj)
}

func TestImportExportBox(t *testing.T) {
	obj, err := Import("test/box.obj", "test/box.mtl", "test/box_material.txt")
	if err != nil {
		t.Fatalf("import: %s", err)
	}
	obj.Name = "box"
	err = Export(obj, "test/tmp.obj", "test/tmp.mtl", "test/tmp_material.txt")
	if err != nil {
		t.Fatalf("export: %s", err)
	}
}
