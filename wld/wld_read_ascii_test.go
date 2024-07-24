package wld

import (
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/common"
)

func TestReadDMSpriteDef2(t *testing.T) {
	if os.Getenv("SINGLE_TEST") != "1" {
		t.Skip("skipping test; SINGLE_TEST not set")
	}
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	tests := []struct {
		asciiName string
		wantErr   bool
	}{}
	if !common.IsTestExtensive() {
		tests = []struct {
			asciiName string
			wantErr   bool
		}{
			//{"all.spk", false},
			{"fis.spk", false},
			//{"pre.spk", false},
		}
	}
	for _, tt := range tests {
		t.Run(tt.asciiName, func(t *testing.T) {

			wld := &Wld{}
			err := wld.ReadAscii(fmt.Sprintf("testdata/%s", tt.asciiName))
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.asciiName, err.Error())
			}
		})
	}

}
