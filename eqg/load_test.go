package eqg

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestEQGLoad(t *testing.T) {

	inFile := "test/box.eqg"
	f, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(inFile)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save(fmt.Sprintf("%s.png", inFile))

	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
}

func TestLoadSaveLoad(t *testing.T) {
	inFile1 := "test/oasis.eqg"
	outFile1 := "test/oasis_out.eqg"
	inFile2 := "test/oasis_out.eqg"

	f, err := os.Open(inFile1)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	d, err := dump.New(inFile1)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	archive, err := New(inFile1)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = archive.Load(f)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	w, err := os.Create(outFile1)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = archive.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err)
	}
	err = d.Save(inFile1 + ".png")
	if err != nil {
		t.Fatalf("dump save: %s", err)
	}
	dump.Close()

	d, err = dump.New(inFile2)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}

	r, err := os.Open(inFile2)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = archive.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	err = d.Save(inFile2 + ".png")
	if err != nil {
		t.Fatalf("dump save2: %s", err)
	}
	dump.Close()
}
