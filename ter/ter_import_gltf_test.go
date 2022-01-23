package ter

import (
	"fmt"
	"os"
	"testing"
)

func TestGLTFImport(t *testing.T) {
	e := &TER{}
	path := "../eq/tmp/ecommons.gltf"
	err := e.ImportGLTF(path)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	w, err := os.Create("../eq/tmp/out.mod")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}

func TestGLTFImportSave(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &TER{}
	path := "../eq/tmp/ecommons.gltf"
	err := e.ImportGLTF(path)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	w, err := os.Create("../eq/tmp/out.ter")
	if err != nil {
		t.Fatalf("create: %s", err)
	}

	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}
