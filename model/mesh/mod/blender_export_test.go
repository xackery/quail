package mod

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestMOD_BlenderExport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name    string
		eqgPath string
		dstDir  string
		wantErr bool
	}{
		{name: "test", eqgPath: fmt.Sprintf("%s/it13900.eqg", eqPath), dstDir: "test/", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseName := filepath.Base(tt.eqgPath)
			baseName = strings.TrimSuffix(baseName, ".eqg")

			pfs, err := eqg.NewFile(tt.eqgPath)
			if err != nil {
				t.Fatalf("failed to open eqg file: %s", err.Error())
			}

			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".mod" {
					continue
				}
				e, err := New(baseName+".mod", pfs)
				if err != nil {
					t.Fatalf("failed to create mod: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("failed to decode mod: %s", err.Error())
				}

				if err := e.BlenderExport(tt.dstDir); (err != nil) != tt.wantErr {
					t.Errorf("mod.BlenderExport() error = %v, wantErr %v", err, tt.wantErr)
				}
				break
			}

		})
	}
}
