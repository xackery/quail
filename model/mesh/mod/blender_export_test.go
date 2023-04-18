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

func TestBlender_Export(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name      string
		modelName string
		wantErr   bool
	}{
		{name: "it13926.eqg", wantErr: false},
		{name: "it12095.eqg", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("Failed to open eqg file: %s", err.Error())
			}

			dir := fmt.Sprintf("test/_%s", tt.name)
			isFound := false
			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".mod" {
					continue
				}
				if tt.modelName != "" && !strings.Contains(fe.Name(), tt.modelName) {
					continue
				}

				isFound = true

				e, err := New(fe.Name(), pfs)
				if err != nil {
					t.Fatalf("Failed to new mod: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("Failed to decode mod: %s", err.Error())
				}

				err = e.BlenderExport(dir)
				if err != nil {
					t.Fatalf("Failed to export mod: %s", err.Error())
				}
			}
			if !isFound {
				t.Fatalf("mod %s not found", tt.name)
			}
		})
	}
}
