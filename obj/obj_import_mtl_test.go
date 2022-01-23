package obj

import (
	"os"
	"testing"
)

func TestImportMtl(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	obj := &ObjData{}
	err := importMtl(obj, "../eq/soldungb/cache/soldungb.mtl")
	if err != nil {
		t.Fatalf("importMtl: %s", err)
	}
	t.Fatalf("%+v", obj)
}
