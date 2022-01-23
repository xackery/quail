package s3d

import (
	"os"
	"testing"
)

func TestWldLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	f, err := os.Open("../eq/tmp/lights.wld")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	e := &Wld{}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
}
