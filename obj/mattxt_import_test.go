package obj

import (
	"testing"
)

func TestMattxtImport(t *testing.T) {
	req := &ObjRequest{
		Data:       &ObjData{},
		MattxtPath: "test/box_material.txt",
	}
	err := mattxtImport(req)
	if err != nil {
		t.Fatalf("mattxtImport: %s", err)
	}
}
