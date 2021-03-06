package ter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/eqg"
)

func TestArchiveExportSample(t *testing.T) {
	tests := []struct {
		category string
	}{
		{category: "box"},
		//{category: "steamfontmts"},
		//{category: "broodlands"},
		//{category: "steppes"},
	}
	for _, tt := range tests {
		isDumpEnabled := false

		eqgFile := fmt.Sprintf("test/%s.eqg", tt.category)
		outFile := fmt.Sprintf("test/%s_out.ter", tt.category)

		var err error

		a, err := eqg.New(tt.category)
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		r, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		err = a.Decode(r)
		if err != nil {
			t.Fatalf("decode: %s", err)
		}

		e, err := New(tt.category, a)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		for _, fileEntry := range a.Files() {
			if filepath.Ext(fileEntry.Name()) != ".ter" {
				continue
			}
			if isDumpEnabled {
				dump.New(tt.category)
				defer dump.WriteFileClose(fmt.Sprintf("test/%s_eqg_%s", tt.category, fileEntry.Name()))
			}

			terBuf := bytes.NewReader(fileEntry.Data())
			err = e.Decode(terBuf)
			if err != nil {
				t.Fatalf("decode %s: %s", fileEntry.Name(), err)
			}

			w, err := os.Create(outFile)
			if err != nil {
				t.Fatalf("create %s", err)
			}
			defer w.Close()

			err = e.Encode(w)
			if err != nil {
				t.Fatalf("encode: %s", err)
			}

			dump.WriteFileClose(fmt.Sprintf("test/%s_eqg_%s", tt.category, fileEntry.Name()))
		}
	}
}
