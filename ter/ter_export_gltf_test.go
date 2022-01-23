package ter

import (
	"os"
	"testing"
)

func TestGLTFExport(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &TER{
		name: "soldungb",
	}
	path := "../eq/tmp/soldungb.ter"
	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()
	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	err = e.ExportGLTF("../eq/tmp/out.glb")
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}
