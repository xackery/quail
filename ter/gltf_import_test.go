package ter

import (
	"fmt"
	"os"
	"testing"

	"github.com/qmuntal/gltf"
)

func TestGLTFImportExportBoxGLTF(t *testing.T) {
	path := "test/box.eqg"
	inFile := "test/box_out.gltf"

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.GLTFImport(inFile)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

}

func TestGLTFBoxVerify(t *testing.T) {
	inFile := "test/box.gltf"
	outFile := "test/box_out.gltf"

	doc, err := gltf.Open(inFile)
	if err != nil {
		t.Fatalf("gltf.Open: %s", err)
	}
	gltf.Save(doc, outFile)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestGLTFBoxSanity(t *testing.T) {
	path := "test/box.eqg"
	inFile := "test/box.gltf"
	outFile := "test/box_out.gltf"

	//validator: https://github.khronos.org/glTF-Validator/

	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	err = e.GLTFImport(inFile)
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	err = e.GLTFExport(w)
	if err != nil {
		t.Fatalf("GLTFExport: %s", err)
	}

	docOriginal, err := gltf.Open(inFile)
	if err != nil {
		t.Fatalf("gltf.Open docOriginal: %s", err)
	}
	if docOriginal == nil {
		t.Fatalf("docOriginal nil")
	}
	docExport, err := gltf.Open(outFile)
	if err != nil {
		t.Fatalf("gltf.Open docExport: %s", err)
	}
	if docExport == nil {
		t.Fatalf("docExport nil")
	}
	/*for i := 0; i < len(docOriginal.Accessors); i++ {
		if docOriginal.Accessors[i].Count != docExport.Accessors[i].Count {
			t.Fatalf("mismatch accessors %d count: wanted %d, got %d", i, docOriginal.Accessors[i].Count, docExport.Accessors[i].Count)
		}
	}*/
	fmt.Printf("dump: %+v\n", e)
}
