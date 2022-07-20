package eqg

import (
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestDecodeEncodeCompareNoDebug(t *testing.T) {
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
	err = e.Decode(f)
	if err != nil {
		t.Fatalf("decode: %s", err)
	}

	w, err := os.Create(out)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = e.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
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
	dump.New(inFile)
	defer dump.WriteFileClose(inFile + ".png")
	archive, err := New(inFile)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = archive.Decode(f)
	if err != nil {
		t.Fatalf("decode: %s", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err)
	}
	defer w.Close()

	err = archive.Encode(w)
	if err != nil {
		t.Fatalf("encode: %s", err)
	}
}
