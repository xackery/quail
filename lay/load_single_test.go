package lay

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
)

func TestLoadSingleTest(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		return
	}
	isDump := true
	tests := []struct {
		category string
	}{
		{category: "steamfontmts"},
		{category: "lth"},
	}
	for _, tt := range tests {

		fmt.Println("loading", tt.category)
		eqgFile := fmt.Sprintf("test/eq/%s.eqg", tt.category)

		ra, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		defer ra.Close()
		a, err := eqg.New(tt.category)
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		err = a.Load(ra)
		if err != nil {
			t.Fatalf("load eqg: %s", err)
		}

		files := a.Files()
		for _, layEntry := range files {
			if filepath.Ext(layEntry.Name()) != ".lay" {
				continue
			}
			fmt.Println(layEntry.Name())

			var d *dump.Dump
			if isDump {
				d, err = dump.New(layEntry.Name())
				if err != nil {
					t.Fatalf("dump.New: %s", err)
				}
			}
			r := bytes.NewReader(layEntry.Data())

			e, err := NewEQG(layEntry.Name(), a)
			if err != nil {
				t.Fatalf("new: %s", err)
			}

			err = e.Load(r)
			if err != nil {
				t.Fatalf("load %s: %s", layEntry.Name(), err)
			}

			if isDump {
				err = d.Save(fmt.Sprintf("test/eq/%s_eqg_%s.png", tt.category, layEntry.Name()))
				if err != nil {
					t.Fatalf("save: %s", err)
				}
			}
			fmt.Println(e.layers)
		}
	}
}
