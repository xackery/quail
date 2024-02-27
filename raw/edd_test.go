package raw

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestEddRead(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	//dirTest := common.DirTest()

	tests := []struct {
		name string
	}{
		{name: "actoremittersnew.edd"}, // FIXME: proper edd read support
		{name: "environmentemittersnew.edd"},
		{name: "spellsnew.edd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := os.Open(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open %s: %s", tt.name, err.Error())
			}
			defer r.Close()
			edd := &Edd{}
			err = edd.Read(r)
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}
		})
	}
}

func TestEddWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	//dirTest := common.DirTest()

	tests := []struct {
		name string
	}{
		{name: "actoremittersnew.edd"}, // FIXME: proper edd write support
		{name: "environmentemittersnew.edd"},
		{name: "spellsnew.edd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := os.Open(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open %s: %s", tt.name, err.Error())
			}
			defer r.Close()
			edd := &Edd{}
			err = edd.Read(r)
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			buf := bytes.NewBuffer(nil)
			err = edd.Write(buf)
			if err != nil {
				t.Fatalf("failed to write %s: %s", tt.name, err.Error())
			}

		})
	}
}
