package wld

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/common"
)

func TestAsciiRead(t *testing.T) {
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
			//{"all/all.spk", false},
			//{"fis/fis.spk", false},
			{"pre/pre.spk", false},
		}
	}
	for _, tt := range tests {
		t.Run(tt.asciiName, func(t *testing.T) {

			wld := &Wld{}
			err := wld.ReadAscii(fmt.Sprintf("testdata/%s", tt.asciiName))
			if err != nil {
				t.Fatalf("Failed readascii: %s", err.Error())
			}
		})
	}
}

func TestAsciiReadWriteRead(t *testing.T) {
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
			//{"all/all.spk", false},
			//{"fis/fis.spk", false},
			{"pre/pre.spk", false},
		}
	}
	for _, tt := range tests {
		t.Run(tt.asciiName, func(t *testing.T) {

			wld := &Wld{
				FileName: filepath.Base(tt.asciiName),
			}
			err := wld.ReadAscii(fmt.Sprintf("testdata/%s", tt.asciiName))
			if err != nil {
				t.Fatalf("Failed readascii: %s", err.Error())
			}

			err = wld.WriteAscii("testdata/temp/", true)
			if err != nil {
				t.Fatalf("Failed writeascii: %s", err.Error())
			}

			wld2 := &Wld{}
			err = wld2.ReadAscii("testdata/temp/" + wld.FileName)
			if err != nil {
				t.Fatalf("Failed readascii: %s", err.Error())
			}

		})
	}
}
