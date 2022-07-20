package prt

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
)

func TestDecode(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "test/eq/djm.eqg"
	inFile := "djm.prt"

	archive, err := eqg.NewFile(path)
	if err != nil {
		t.Fatalf("eqg new: %s", err)
	}

	dump.New(inFile)
	defer dump.WriteFileClose(fmt.Sprintf("test/eq/%s.png", inFile))

	e, err := New(inFile, archive)
	if err != nil {
		t.Fatalf("prt new: %s", err)
	}
	data, err := archive.File(inFile)
	if err != nil {
		t.Fatalf("file: %s", err)
	}

	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode prt: %s", err)
	}
}
