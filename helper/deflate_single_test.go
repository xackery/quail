package helper

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func TestDeflate(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	data, err := os.ReadFile("../scripts/deflate/data")
	if err != nil {
		t.Fatalf("readfile: %s", err.Error())
	}
	out, err := Deflate(data)
	if err != nil {
		t.Fatalf("deflate: %s", err.Error())
	}
	//fmt.Println(hex.Dump(out[0:16]))
	fmt.Println(hex.Dump(out))
	//00000000  26 01 00 00 32 04 00 00  58 85 75 53 db 4a 03 31  |&...2...X.uS.J.1|
	//0000  40 1 0 0 32 4 0 0 78 5e 75 53 ffffffcb 4a 43 31 ||
}
