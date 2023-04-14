package lit

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestLIT_Decode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name    string
		eqgPath string
		litFile string
		dstDir  string
		wantErr bool
	}{

		{name: "test", eqgPath: fmt.Sprintf("%s/steamfontmts.eqg", eqPath), litFile: "akanon_entry_obj_lavalightt01.lit", dstDir: "test", wantErr: false},
		{name: "test", eqgPath: fmt.Sprintf("%s/steamfontmts.eqg", eqPath), litFile: "akanon_entry_obj_pipe_outera09.lit", dstDir: "test", wantErr: false},
	}
	for _, tt := range tests {
		baseName := filepath.Base(tt.eqgPath)
		baseName = strings.TrimSuffix(baseName, ".eqg")

		pfs, err := eqg.NewFile(tt.eqgPath)
		if err != nil {
			t.Fatalf("failed to open eqg file: %s", err.Error())
		}

		for _, fe := range pfs.Files() {
			if filepath.Ext(fe.Name()) != ".lit" {
				continue
			}
			if fe.Name() != tt.litFile {
				continue
			}
			e, err := New(baseName+".lit", pfs)
			if err != nil {
				t.Fatalf("failed to create lit: %s", err.Error())
			}

			err = e.Decode(bytes.NewReader(fe.Data()))
			if err != nil {
				t.Fatalf("failed to decode lit: %s", err.Error())
			}

			if err := e.BlenderExport(tt.dstDir); (err != nil) != tt.wantErr {
				t.Errorf("lit.BlenderExport() error = %v, wantErr %v", err, tt.wantErr)
			}
			break
		}
	}
}
