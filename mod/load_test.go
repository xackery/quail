package mod

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	path := "test/chair.mod"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save("test/out.png")

	e := &MOD{}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

}
