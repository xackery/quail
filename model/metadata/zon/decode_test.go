package zon

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestZON_Decode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		//{name: "bloodfields.eqg", wantErr: false}, //v1
		{name: "anguish.eqg", wantErr: false}, //v1
		{name: "beastdomain.eqg", wantErr: false},
		{name: "steamfontmts.eqg", wantErr: false}, //v4
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("Failed to open eqg file: %s", err.Error())
			}

			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".zon" {
					continue
				}
				e, err := New(fe.Name(), pfs)
				if err != nil {
					t.Fatalf("Failed to new zon: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("Failed to decode zon: %s", err.Error())
				}
				break
			}

		})
	}
}
