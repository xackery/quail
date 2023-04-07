package ter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestTER_BlenderExport(t *testing.T) {
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
		{name: "test", eqgPath: fmt.Sprintf("%s/arena.eqg", eqPath), dstDir: "test/", wantErr: false},
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
				if filepath.Ext(fe.Name()) != ".ter" {
					continue
				}
				e, err := New(baseName+".ter", pfs)
				if err != nil {
					t.Fatalf("failed to create ter: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("failed to decode ter: %s", err.Error())
				}

				if err := e.BlenderExport(tt.dstDir); (err != nil) != tt.wantErr {
					t.Errorf("TER.BlenderExport() error = %v, wantErr %v", err, tt.wantErr)
				}
				break
			}

		})
	}
}
