package eqg

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/xackery/quail/dump"
)

func TestSaveBoxEQG(t *testing.T) {
	inFile := "test/box.eqg"
	outFile := "test/tmp.eqg"
	isDump := true

	if isDump {
		d, err := dump.New(inFile)
		if err != nil {
			t.Fatalf("dump.New: %s", err)
		}
		defer d.Save(fmt.Sprintf("%s.png", inFile))
	}
	r, err := os.Open(inFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer r.Close()

	e, err := New(inFile)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	err = e.Load(r)
	if err != nil {
		t.Fatalf("load: %s", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create: %s", err.Error())
	}
	defer w.Close()
	err = e.Save(w)
	if err != nil {
		t.Fatalf("save: %s", err.Error())
	}
	compareFile(t, inFile, outFile)
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
