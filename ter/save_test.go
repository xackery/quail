package ter

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	var err error
	e := &TER{}
	err = e.AddMaterial("test", "test2")
	if err != nil {
		t.Fatalf("addModel: %s", err)
	}
	err = e.AddMaterialProperty("test", "testProp", 0, 1, 0)
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
