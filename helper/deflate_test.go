package helper

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDeflate(t *testing.T) {
	out, err := Deflate([]byte("test"))
	if err != nil {
		t.Fatalf("deflate: %s", err.Error())
	}
	fmt.Println(hex.Dump(out))
}
