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
		//{category: "steamfontmts"},
		//{category: "lth"},
		//{category: "bxi"}, //bixie, two models and multiple layers
		{category: "alg"},
		//{category: "cak"},
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

			dump.New(layEntry.Name())

			if isDump {
				dump.New(layEntry.Name())
			}
			defer dump.WriteFileClose(fmt.Sprintf("test/eq/%s_%s", tt.category, layEntry.Name()))
			r := bytes.NewReader(layEntry.Data())

			e, err := New(layEntry.Name(), a)
			if err != nil {
				t.Fatalf("new: %s", err)
			}

			err = e.Load(r)
			if err != nil {
				t.Fatalf("load %s: %s", layEntry.Name(), err)
			}
			dump.WriteFileClose(fmt.Sprintf("test/eq/%s_%s", tt.category, layEntry.Name()))
		}
	}
}
