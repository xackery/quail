package raw

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestEcoWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()

	tests := []struct {
		name    string
		ecoName string
	}{
		{name: "arcstone.eqg", ecoName: "farstone_base.eco"},

		//{name: "ggy.eqg", ecoName: "ggy.eco"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.name))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".eco" {
					continue
				}
				eco := &Eco{}

				err = eco.Read(bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = eco.Write(buf)
				if err != nil {
					t.Fatalf("failed to write %s: %s", tt.name, err.Error())
				}

				eco2 := &Eco{}
				err = eco2.Read(bytes.NewReader(buf.Bytes()))
				if err != nil {
					t.Fatalf("failed to read %s: %s", tt.name, err.Error())
				}

				err = helper.ByteCompareTest(file.Data(), buf.Bytes())
				if err != nil {
					t.Fatalf("%s byteCompare: %s", tt.name, err)
				}
			}
		})
	}
}
