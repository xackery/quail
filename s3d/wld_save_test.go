package s3d

import (
	"bytes"
	"os"
	"testing"
)

func TestWldSave(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	f := &bytes.Buffer{}
	e := &Wld{}
	err := e.Save(f)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
}
