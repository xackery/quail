package mod

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/gltf"
)

func TestGLTFDecode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	filePath := "test/"
	inFile := "test/ecommons.gltf"
	outFile := "test/ecommons_gltfimport.mod"

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	gdoc, err := gltf.Open(inFile)
	if err != nil {
		t.Fatalf("gltf open: %s", err)
	}
	err = e.GLTFDecode(gdoc)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}

func TestGLTFDecodeEncode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	filePath := "test/"
	inFile := "test/ecommons.gltf"
	outFile := "test/ecommons_gltfimportsave.mod"

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
	}
	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	gdoc, err := gltf.Open(inFile)
	if err != nil {
		t.Fatalf("gltf open: %s", err)
	}
	err = e.GLTFDecode(gdoc)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}
	fmt.Printf("dump: %+v\n", e)
}
