package ter

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/_steamfontmts.eqg/steamfontmts.ter"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save("../eq/tmp/out.png")

	e := &TER{}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
}

func TestLoadSaveLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/_ecommons.eqg/ecommons.ter"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	e := &TER{}
	err = e.Load(f)
	if err != nil {
		d.Save("../eq/tmp/out.png")
		t.Fatalf("load: %s", err)
	}
	w, err := os.Create("../eq/tmp/out.ter")
	if err != nil {
		d.Save("../eq/tmp/out.png")
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = e.Save(w)
	if err != nil {
		d.Save("../eq/tmp/out.png")
		t.Fatalf("save: %s", err)
	}
	d.Save("../eq/tmp/out.png")
	dump.Close()

	path = "../eq/tmp/out.ter"
	d, err = dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save("../eq/tmp/out2.png")
	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("reload: %s", err)
	}

}
