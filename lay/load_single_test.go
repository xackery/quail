package lay

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
	path := "test/eq/lth.eqg"
	inFile := "lth.lay"

	a, err := eqg.New(path)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	ra, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer ra.Close()
	err = a.Load(ra)
	if err != nil {
		t.Fatalf("archive load: %s", err)
	}

	d, err := dump.New(inFile)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save(fmt.Sprintf("test/eq/%s.png", inFile))

	e, err := NewEQG(inFile, a)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	data, err := a.File(inFile)
	if err != nil {
		t.Fatalf("file: %s", err)
	}

	err = e.Load(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("load mds: %s", err)
	}
}
