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
