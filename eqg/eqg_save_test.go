package eqg

import (
	"os"
	"testing"
)

func TestSave(t *testing.T) {

	e := &EQG{}
	err := e.Add("test.txt", []byte("test"))
	if err != nil {
		t.Fatalf("add: %s", err.Error())
	}
	f, err := os.Create("test/out.eqg")
	if err != nil {
		t.Fatalf("create: %s", err.Error())
	}
	err = e.Save(f)
	if err != nil {
		t.Fatalf("save: %s", err.Error())
	}
	f.Close()
}
