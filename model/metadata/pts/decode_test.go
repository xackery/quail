package pts

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/pfs/eqg"
)

func TestPTS_Decode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "sin.eqg", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pfs, err := eqg.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("Failed to open eqg file: %s", err.Error())
			}

			isFound := false
			for _, fe := range pfs.Files() {
				if filepath.Ext(fe.Name()) != ".pts" {
					continue
				}
				isFound = true
				e, err := New(fe.Name(), pfs)
				if err != nil {
					t.Fatalf("Failed to new pts: %s", err.Error())
				}

				err = e.Decode(bytes.NewReader(fe.Data()))
				if err != nil {
					t.Fatalf("Failed to decode pts: %s", err.Error())
				}
				break
			}
			if !isFound {
				t.Fatalf("pts %s not found", tt.name)
			}

		})
	}
}
