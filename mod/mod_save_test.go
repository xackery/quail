package mod

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	var err error

	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.AddMaterial("test", "test2")
	if err != nil {
		t.Fatalf("addModel: %s", err)
	}
	err = e.AddMaterialProperty("test", "testProp", 0, "1")
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
