package wld

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xackery/quail/pfs/s3d"
)

func TestWLD_BlenderExport(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name    string
		dstDir  string
		wantErr bool
	}{
		{name: "ggl_chr.s3d", dstDir: "test/", wantErr: false},
		{name: "arena.s3d", dstDir: "test/", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("%s/%s", eqPath, tt.name)
			baseName := filepath.Base(path)
			baseName = strings.TrimSuffix(baseName, ".s3d")

			pfs, err := s3d.NewFile(path)
			if err != nil {
				t.Fatalf("Failed to open s3d file: %s", err.Error())
			}

			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".wld" {
					continue
				}
				e, err := New(baseName+".wld", pfs)
				if err != nil {
					t.Fatalf("Failed to create wld: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("Failed to decode wld: %s", err.Error())
				}

				if err := e.BlenderExport(tt.dstDir); (err != nil) != tt.wantErr {
					t.Errorf("wld.BlenderExport() error = %v, wantErr %v", err, tt.wantErr)
				}
				break
			}

		})
	}
}
