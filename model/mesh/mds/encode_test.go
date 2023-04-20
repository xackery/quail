package mds

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestMDS_Encode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name      string
		mdselName string
		version   uint32
		wantErr   bool
	}{
		{name: "it12095.eqg", mdselName: "gnome_wave.mds", version: 3, wantErr: false},
		//{name: "sin.eqg", mdselName: "sin.mds", version: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("Failed to open eqg file: %s", err.Error())
			}

			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".mds" {
					continue
				}
				if tt.mdselName != "" && tt.mdselName != fe.Name() {
					continue
				}
				e, err := New(fe.Name(), pfs)
				if err != nil {
					t.Fatalf("Failed to new mds: %s", err.Error())
				}
				e.version = tt.version

				outDir := fmt.Sprintf("test/_%s/test_data/", tt.name)
				err = os.MkdirAll(outDir, 0755)
				if err != nil {
					t.Fatalf("Failed to create dir: %s", err.Error())
				}

				err = os.WriteFile(fmt.Sprintf("%s/%s-raw.mds", outDir, fe.Name()), fe.Data(), 0755)
				if err != nil {
					t.Fatalf("Failed to write file: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("Failed to decode mds: %s", err.Error())
				}

				w, err := os.Create(fmt.Sprintf("%s/%s-encoded.mds", outDir, fe.Name()))
				if err != nil {
					t.Fatalf("Failed to create file: %s", err.Error())
				}
				defer w.Close()

				err = e.Encode(w)
				if err != nil {
					t.Fatalf("Failed to encode mds: %s", err.Error())
				}
			}

		})
	}
}
