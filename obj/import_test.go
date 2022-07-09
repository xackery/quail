package obj

import "testing"

func TestImportBox(t *testing.T) {
	obj, err := Import("test/box.obj", "test/box.mtl", "test/box_material.txt")
	if err != nil {
		t.Fatalf("import: %s", err)
	}
	if obj == nil {
		t.Fatalf("empty object")
	}
}
