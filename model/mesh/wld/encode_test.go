package wld

import (
	"bytes"
	"os"
	"testing"

	"github.com/xackery/quail/pfs/s3d"
)

func TestWldEncode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	archive, err := s3d.New("test")
	if err != nil {
		t.Fatalf("s3d new: %s", err)
	}
	f := &bytes.Buffer{}
	e, err := New("out", archive)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Encode(f)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
}
