package s3d

import (
	"bytes"
	"testing"
)

func TestWldSave(t *testing.T) {
	f := &bytes.Buffer{}
	e := &Wld{}
	err := e.Save(f)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
}
