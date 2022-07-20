package mds

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/eqg"
)

func TestSave(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	var err error

	filePath := "test/"
	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("newPath: %s", err)
	}
	e, err := New("out", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.MaterialAdd("test", "test2")
	if err != nil {
		t.Fatalf("addModel: %s", err)
	}
	err = e.MaterialPropertyAdd("test", "testProp", 0, "1")
	if err != nil {
		t.Fatalf("addMaterialProperty: %s", err)
	}
	buf := bytes.NewBuffer(nil)

	err = e.Save(buf)
	if err != nil {
		t.Fatalf("save: %s", err.Error())
	}
	fmt.Println(hex.Dump(buf.Bytes()))
}

func TestSaveEQG(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	var err error

	path, err := common.NewPath("test/")
	if err != nil {
		t.Fatalf("newPath: %s", err)
	}

	archive, err := eqg.New("gbm")
	if err != nil {
		t.Fatalf("new: %s", err)
	}

	e, err := New("gbm", path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.GLTFImport("test/box.gltf")
	if err != nil {
		t.Fatalf("import %s: %s", path, err)
	}

	err = e.ArchiveExport(archive)
	if err != nil {
		t.Fatalf("archive export: %s", err)
	}

	w, err := os.Create("test/gbm.eqg")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = archive.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}

}
