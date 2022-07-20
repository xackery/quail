package eqg

import (
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

	dump.New(inFile)
	defer dump.WriteFileClose(inFile)

	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Decode(f)
	if err != nil {
		t.Fatalf("decode: %s", err)
	}
}

func TestDecodeEncodeDecode(t *testing.T) {
	inFile1 := "test/oasis.eqg"
	outFile1 := "test/oasis_out.eqg"
	inFile2 := "test/oasis_out.eqg"

	f, err := os.Open(inFile1)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()

	dump.New(inFile1)
	defer dump.WriteFileClose(inFile1)

	archive, err := New(inFile1)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = archive.Decode(f)
	if err != nil {
		t.Fatalf("decode: %s", err)
	}

	w, err := os.Create(outFile1)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()
	err = archive.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}

	dump.WriteFileClose(inFile1)
	dump.New(inFile2)
	defer dump.WriteFileClose(inFile2)

	r, err := os.Open(inFile2)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = archive.Decode(r)
	if err != nil {
		t.Fatalf("decode: %s", err)
	}
}
