package eqg

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/zon"
)

func TestArchiveExportSample(t *testing.T) {
	tests := []struct {
		category    string
		outCategory string
	}{
		{category: "box", outCategory: "arena"},
		//{category: "steamfontmts"},
		//{category: "broodlands"},
		//{category: "steppes"},
	}
	for _, tt := range tests {

		eqgFile := fmt.Sprintf("test/%s.eqg", tt.category)
		outFile := fmt.Sprintf("test/%s_out.eqg", tt.outCategory)

		var err error

		archive, err := New(tt.category)
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		r, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		err = archive.Load(r)
		if err != nil {
			t.Fatalf("load: %s", err)
		}

		outArchive, err := New(tt.category)
		if err != nil {
			t.Fatalf("new out eqg: %s", err)
		}

		e, err := zon.New(tt.outCategory, archive)
		if err != nil {
			t.Fatalf("new: %s", err)
		}

		data, err := archive.File(fmt.Sprintf("%s.zon", tt.category))
		if err != nil {
			t.Fatalf("archive.file: %s", err)
		}
		err = e.Load(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("load: %s", err)
		}

		err = e.ArchiveExport(outArchive)
		if err != nil {
			t.Fatalf("archive export: %s", err)
		}

		w, err := os.Create(outFile)
		if err != nil {
			t.Fatalf("create: %s", err)
		}
		defer w.Close()
		err = outArchive.Save(w)
		if err != nil {
			t.Fatalf("save: %s", err)
		}
	}
}
