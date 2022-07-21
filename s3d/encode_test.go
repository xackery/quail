package s3d

import (
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
		{category: "crushbone", outCategory: "arena"},
		//{category: "steamfontmts"},
		//{category: "broodlands"},
		//{category: "steppes"},
	}
	for _, tt := range tests {

		eqgFile := fmt.Sprintf("test/eq/%s.s3d", tt.category)
		outFile := fmt.Sprintf("test/eq/%s_out.s3d", tt.outCategory)

		var err error

		archive, err := New(tt.category)
		if err != nil {
			t.Fatalf("eqg.New: %s", err)
		}
		r, err := os.Open(eqgFile)
		if err != nil {
			t.Fatalf("%s", err)
		}
		err = archive.Decode(r)
		if err != nil {
			t.Fatalf("decode: %s", err)
		}

		outArchive, err := New(tt.category)
		if err != nil {
			t.Fatalf("new out eqg: %s", err)
		}

		e, err := zon.NewFile(tt.outCategory, archive, fmt.Sprintf("%s.zon", tt.category))
		if err != nil {
			t.Fatalf("new: %s", err)
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
		err = outArchive.Encode(w)
		if err != nil {
			t.Fatalf("encode: %s", err)
		}
	}
}
