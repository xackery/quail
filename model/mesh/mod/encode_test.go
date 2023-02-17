package mod

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
)

func TestEncode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	var err error

	filePath := "test/"

	path, err := common.NewPath(filePath)
	if err != nil {
		t.Fatalf("path: %s", err)
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

	err = e.Encode(buf)
	if err != nil {
		t.Fatalf("encode: %s", err.Error())
	}
	fmt.Println(hex.Dump(buf.Bytes()))
}
