package ter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestTER_Encode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "arena.eqg", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("Failed to open eqg file: %s", err.Error())
			}

			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".ter" {
					continue
				}
				e, err := New(fe.Name(), pfs)
				if err != nil {
					t.Fatalf("Failed to new ter: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("Failed to decode ter: %s", err.Error())
				}

				out := fmt.Sprintf("test/_%s/%s", tt.name, fe.Name())
				err = os.MkdirAll(filepath.Dir(out), 0755)
				if err != nil {
					t.Fatalf("Failed to create dir: %s", err.Error())
				}
				w, err := os.Create(out)
				if err != nil {
					t.Fatalf("Failed to create file: %s", err.Error())
				}
				defer w.Close()

				err = e.Encode(w)
				if err != nil {
					t.Fatalf("Failed to encode ter: %s", err.Error())
				}
				break
			}

		})
	}
}
