package wld

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/_crushbone.s3d/crushbone.wld"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.New: %s", err)
	}
	defer d.Save("../eq/tmp/out.png")
	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(f)
	if err != nil {
		d.Save("../eq/tmp/out.png")
		t.Fatalf("load: %s", err)
	}
}
