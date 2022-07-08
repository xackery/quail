package eqg

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/arena.eqg"
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

func TestLoadSaveCompareNoDebug(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	in := "/Users/xackery/Documents/games/EverQuest.app/drive_c/rebuildeq/arena.eqg"
	out := "/Users/xackery/Documents/games/EverTest.app/drive_c/rebuildeq/arena.eqg"
	f, err := os.Open(in)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	e, err := New(in)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	w, err := os.Create(out)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
}

func TestLoadSaveLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/arena.eqg"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	w, err := os.Create("../eq/tmp/arena.eqg")
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	d.Save("../eq/tmp/arena_original.png")
	dump.Close()

	path = "../eq/tmp/arena.eqg"
	d, err = dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	d.Save("../eq/tmp/arena_new.png")

}
