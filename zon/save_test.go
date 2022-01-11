package zon

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/g3n/engine/math32"
)

func TestSave(t *testing.T) {
	var err error
	z := &ZON{}
	err = z.AddModel("test")
	if err != nil {
		t.Fatalf("addModel: %s", err)
	}
	err = z.AddObject("test", "test01", math32.Vector3{X: 1, Y: 2, Z: 3}, math32.Vector3{}, 0)
	if err != nil {
		t.Fatalf("addObject: %s", err)
	}
	buf := bytes.NewBuffer(nil)

	err = z.Save(buf)
	if err != nil {
		t.Fatalf("save: %s", err.Error())
	}
	fmt.Println(hex.Dump(buf.Bytes()))
}
