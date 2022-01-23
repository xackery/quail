package zon

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/_ecommons.eqg/ecommons.zon"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save("../eq/out.png")

	e := &ZON{}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
}

func TestLoadSaveLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/_ecommons.eqg/ecommons.zon"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	e := &ZON{}
	err = e.Load(f)
	if err != nil {
		d.Save("../eq/tmp/out.png")
		t.Fatalf("load: %s", err)
	}
	w, err := os.Create("../eq/tmp/out.zon")
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

	path = "../eq/tmp/out.zon"
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
