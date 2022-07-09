package obj

import (
	"os"
	"testing"
)

func TestMattxtImport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	obj := &ObjData{}
	err := mattxtImport(obj, "../eq/soldungb/cache/soldungb_material.txt")
	if err != nil {
		t.Fatalf("importMatTxt: %s", err)
	}
	t.Fatalf("%+v", obj)
}
