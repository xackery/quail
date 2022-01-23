package obj

import (
	"os"
	"testing"
)

func TestImportObjFile(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	obj := &ObjData{}
	err := importObjFile(obj, "../eq/soldungb/cache/soldungb.obj")
	if err != nil {
		t.Fatalf("importObj: %s", err)
	}
	t.Fatalf("%+v", obj)
}
