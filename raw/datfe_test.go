package raw

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestDatFeWrite(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := helper.DirTest()

	tests := []struct {
		name         string
		pfsName      string
		quadsPerTile int
		wantErr      bool
		isDump       bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.pfsName))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			data, err := pfs.File(tt.name)
			if err != nil {
				t.Fatalf("failed to open %s: %s", tt.name, err.Error())
			}

			dat := &DatIw{}

			err = dat.Read(bytes.NewReader(data))
			if tt.isDump {
				os.WriteFile(fmt.Sprintf("%s/%s.src.dat", dirTest, tt.name), data, 0644)
			}
			if err != nil {
				t.Fatalf("failed to read %s: %s", tt.name, err.Error())
			}

			buf := bytes.NewBuffer(nil)

			err = dat.Write(buf)
			if tt.isDump {
				os.WriteFile(fmt.Sprintf("%s/%s.dst.dat", dirTest, tt.name), buf.Bytes(), 0644)
			}
			if err != nil {
				t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
			}

			dat2 := &DatIw{}
			err = dat2.Read(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("failed to decode2 %s: %s", tt.name, err.Error())
			}

			fmt.Println(tt.pfsName, tt.name, "PASS")
		})
	}
}
