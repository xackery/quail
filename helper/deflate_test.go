package helper

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDeflate(t *testing.T) {
	data, err := ioutil.ReadFile("test/test.txt")
	if err != nil {
		t.Fatalf("readfile: %s", err.Error())
	}
	out, err := Deflate(data)
	if err != nil {
		t.Fatalf("deflate: %s", err.Error())
	}
	fmt.Println(hex.Dump(out))
}
