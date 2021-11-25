package s3d

import (
	"os"
	"testing"
)

func TestWldLoad(t *testing.T) {
	f, err := os.Open("test/clz.wld")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	e := &Wld{}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
}
