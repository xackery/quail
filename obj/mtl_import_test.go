package obj

import (
	"testing"
)

func TestMtlImport(t *testing.T) {
	req := &ObjRequest{
		Obj:     &ObjData{},
		MtlPath: "test/tmp.mtl",
	}
	err := mtlImport(req)
	if err != nil {
		t.Fatalf("importMtl: %s", err)
	}
}
