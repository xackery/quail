package ter

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	//path := "test/ecommons.ter"
	path := "test/soldungb.ter"
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

	e := &TER{}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
}

func TestLoadSaveLoad(t *testing.T) {
	path := "test/soldungb.ter"
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
		d.Save("test/out.png")
		t.Fatalf("load: %s", err)
	}
	w, err := os.Create("test/out.ter")
	if err != nil {
		d.Save("test/out.png")
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = e.Save(w)
	if err != nil {
		d.Save("test/out.png")
		t.Fatalf("save: %s", err)
	}
	d.Save("test/out.png")
	dump.Close()

	path = "test/out.ter"
	d, err = dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save("test/out2.png")
	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("reload: %s", err)
	}

}
