package wld

import (
	"bytes"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/s3d"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/eq/crushbone.s3d"

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	defer f.Close()
	a, err := s3d.New("crushbone.s3d")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = a.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.New: %s", err)
	}
	defer d.Save("test/eq/crushbone.wld.png")
	e, err := NewS3D("out", a)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	data, err := a.File("crushbone.wld")
	if err != nil {
		t.Fatalf("file: %s", err)
	}
	r := bytes.NewReader(data)
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load wld: %s", err)
	}
}
