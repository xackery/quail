package wld

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs/s3d"
)

func TestWLD_Encode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}

	log.SetLogLevel(1)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "gequip.s3d", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := s3d.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("s3d.NewFile: %v", err)
				return
			}

			for _, fe := range pfs.Files() {
				if !strings.HasSuffix(fe.Name(), ".wld") {
					continue
				}

				e, err := New(tt.name, pfs)
				if err != nil {
					t.Fatalf("New: %v", err)
					return
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil && !tt.wantErr {
					t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				w, err := os.Create(fmt.Sprintf("test/%s", fe.Name()))
				if err != nil {
					t.Fatalf("os.Create: %v", err)
					return
				}
				err = e.Encode(w)
				if err != nil {
					t.Fatalf("Encode: %v", err)
					return
				}
				w.Close()
			}

		})
	}
}
