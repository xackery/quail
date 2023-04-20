package mod

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestMOD_Encode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name      string
		modelName string
		version   uint32
		wantErr   bool
	}{
		{name: "it13926.eqg", modelName: "it13926.mod", version: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("Failed to open eqg file: %s", err.Error())
			}

			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".mod" {
					continue
				}
				e, err := New(fe.Name(), pfs)
				if err != nil {
					t.Fatalf("Failed to new mod: %s", err.Error())
				}
				e.version = tt.version

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("Failed to decode mod: %s", err.Error())
				}

				outDir := fmt.Sprintf("test/_%s/test_data/", tt.name)
				err = os.MkdirAll(outDir, 0755)
				if err != nil {
					t.Fatalf("Failed to create dir: %s", err.Error())
				}

				err = os.WriteFile(fmt.Sprintf("%s/%s-raw.mod", outDir, fe.Name()), fe.Data(), 0755)
				if err != nil {
					t.Fatalf("Failed to write file: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("Failed to decode mod: %s", err.Error())
				}

				w, err := os.Create(fmt.Sprintf("%s/%s-encoded.mod", outDir, fe.Name()))
				if err != nil {
					t.Fatalf("Failed to create file: %s", err.Error())
				}
				defer w.Close()

				err = e.Encode(w)
				if err != nil {
					t.Fatalf("Failed to encode mod: %s", err.Error())
				}
			}

		})
	}
}
