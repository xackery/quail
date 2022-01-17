package mod

import (
	"fmt"
	"os"
	"testing"
)

func TestGLTFImport(t *testing.T) {

	e := &MOD{}
	path := "test/cube.gltf"
	err := e.ImportGLTF(path)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	w, err := os.Create("test/out.mod")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}
