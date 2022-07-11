package eqg

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestSave(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}

	outFile := "test/tmp.eqg"
	d, err := dump.New(outFile)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save(fmt.Sprintf("%s.png", outFile))

	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Add("test.txt", []byte("test"))
	if err != nil {
		t.Fatalf("add: %s", err.Error())
	}
	f, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err.Error())
	}
	err = e.Save(f)
	if err != nil {
		t.Fatalf("save: %s", err.Error())
	}
	f.Close()
}

func TestSave2(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e, err := New("out")
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Add("test.txt", []byte("test2"))
	if err != nil {
		t.Fatalf("add: %s", err.Error())
	}
	f, err := os.Create("../eq/tmp/out2.eqg")
	if err != nil {
		t.Fatalf("create: %s", err.Error())
	}
	err = e.Save(f)
	if err != nil {
		t.Fatalf("save: %s", err.Error())
	}
	f.Close()
	compareFile(t, "../eq/tmp/out2.eqg", "../eq/tmp/eqzip-test2.eqg")
}
