package prt

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
)

func TestLoad(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/eq/djm.eqg"
	inFile := "djm.prt"

	archive, err := eqg.New(path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	ra, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer ra.Close()
	err = archive.Load(ra)
	if err != nil {
		t.Fatalf("archive load: %s", err)
	}

	d, err := dump.New(inFile)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save(fmt.Sprintf("test/eq/%s.png", inFile))

	e, err := New(inFile, archive)
	if err != nil {
		t.Fatalf("prt new: %s", err)
	}
	data, err := archive.File(inFile)
	if err != nil {
		t.Fatalf("file: %s", err)
	}

	err = e.Load(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("load prt: %s", err)
	}
}
