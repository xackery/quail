package eqg

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestDump(t *testing.T) {
	f, err := os.Open("test/eqzip-test.eqg")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New()
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	e := &EQG{}
	err = e.Load(f)
	if err != nil {
		d.Save("test/out.png")
		t.Fatalf("load: %s", err)
	}

	err = d.Save("test/out.png")
	if err != nil {
		t.Fatalf("save: %s", err)
	}

}
