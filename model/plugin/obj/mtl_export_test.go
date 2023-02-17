package obj

import (
	"testing"
)

func TestExportMtl(t *testing.T) {
	req := &ObjRequest{
		Data:    &ObjData{},
		MtlPath: "test/tmp.mtl",
	}
	err := mtlExport(req)
	if err != nil {
		t.Fatalf("exportMtl: %s", err)
	}
}
