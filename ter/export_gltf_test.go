package ter

import (
	"os"
	"testing"
)

func TestGLTFExport(t *testing.T) {
	e := &TER{
		name: "soldungb",
	}
	path := "test/soldungb.ter"
	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %s", path, err)
	}
	defer r.Close()
	err = e.Load(r)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	err = e.ExportGLTF("test/out.glb")
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}
