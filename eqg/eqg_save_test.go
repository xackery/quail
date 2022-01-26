package eqg

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestSave(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	path := "../eq/ecommons.eqg"
	d, err := dump.New(path)
	if err != nil {
		t.Fatalf("dump.new: %s", err)
	}
	defer d.Save("../eq/tmp/out.png")

	e := &EQG{}
	err = e.Add("test.txt", []byte("test"))
	if err != nil {
		t.Fatalf("add: %s", err.Error())
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %s", err.Error())
	}
	err = e.Save(f)
	if err != nil {
		t.Fatalf("save: %s", err.Error())
	}
	f.Close()

	r, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %s", err)
	}
	err = e.Load(r, "test")
	if err != nil {
		t.Fatalf("load: %s", err)
	}
	compareFile(t, "../eq/tmp/out.eqg", "../eq/tmp/eqzip-test.eqg")
}

func TestSave2(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	e := &EQG{}
	err := e.Add("test.txt", []byte("test2"))
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

func compareFile(t *testing.T, path1 string, path2 string) {

	f1, err := os.Open(path1)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f1.Close()
	f2, err := os.Open(path2)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f2.Close()
	offset := 0

	fails := 0
	f1Data := []byte{}
	f2Data := []byte{}
	for {

		buf1 := make([]byte, 1)
		buf2 := make([]byte, 1)
		_, err1 := f1.Read(buf1)

		if err1 != nil {
			if err1 == io.EOF {
				break
			}
			buf1[0] = 0
		}

		_, err2 := f2.Read(buf2)
		if err2 != nil {
			if err2 == io.EOF {
				break
			}
			buf2[0] = 0
		}
		f1Data = append(f1Data, buf1[0])
		f2Data = append(f2Data, buf2[0])

		if offset == 0 {
			offset++
			continue
		}
		if buf1[0] != buf2[0] {
			fmt.Println(path1, "\n", hex.Dump([]byte(f1Data)))
			fmt.Println(path2, "\n", hex.Dump([]byte(f2Data)))
			if fails > 0 {
				t.Fatalf("mismatched at position %d (0x%02x) %s has value %d, wanted %d", offset, offset, path1, buf1[0], buf2[0])
			}
			fails++
		}

		offset++
	}
}
