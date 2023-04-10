package mod

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestMOD_Decode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "it13926.eqg", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg file: %s", err.Error())
			}

			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".mod" {
					continue
				}
				e, err := New(fe.Name(), pfs)
				if err != nil {
					t.Fatalf("failed to new mod: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("failed to decode mod: %s", err.Error())
				}
				break
			}

		})
	}
}
