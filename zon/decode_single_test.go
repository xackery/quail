package zon

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
)

func TestDecodeSingle(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	category := "oldbloodfield"
	path := "test/eq/" + category + ".eqg"
	archive, err := eqg.NewFile(path)
	if err != nil {
		t.Fatalf("eqg new: %s", err)
	}

	data, err := archive.File(category + ".zon")
	if err != nil {
		t.Fatalf("eqg.file: %s", err)
	}

	dump.New(path)
	defer dump.WriteFileClose(path)

	e, err := New("out", archive)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode: %s", err)
	}
	fmt.Println(e.ModelNames())
}
