package s3d

import (
	"os"
	"testing"
)

func TestS3DEncode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	f, err := os.Create("../eq/tmp/save.s3d")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	err = e.Encode(f)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

}
