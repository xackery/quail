package zon

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
)

func TestLoadBox(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	category := "oldbloodfield"
	path := "test/eq/" + category + ".eqg"
	archive, err := eqg.New(path)
	if err != nil {
		t.Fatalf("eqg.New: %s", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()
	err = archive.Load(f)
	if err != nil {
		t.Fatalf("eqg.load: %s", err)
	}

	data, err := archive.File(category + ".zon")
	if err != nil {
		t.Fatalf("eqg.file: %s", err)
	}

	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save(fmt.Sprintf("%s.png", path))

	e, err := NewEQG("out", archive)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	fmt.Println(e.ModelNames())
}
