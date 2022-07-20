package eqg

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestLoadSaveCompareNoDebug(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	in := "/Users/xackery/Documents/games/EverQuest.app/drive_c/rebuildeq/arena.eqg"
	out := "/Users/xackery/Documents/games/EverTest.app/drive_c/rebuildeq/arena.eqg"
	//in := "test/oasis.eqg"
	//out := "/Users/xackery/Documents/games/EverTest.app/drive_c/rebuildeq/oasis.eqg"
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

func TestLoadSampleCompare(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	inFile := "test/oasis.eqg"
	outFile := "test/oasis_out.eqg"

	f, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(inFile)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	archive, err := New(inFile)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = archive.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = archive.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	d.Save(inFile + ".png")
}
