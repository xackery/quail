package raw

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestDatWtrWrite(t *testing.T) {
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
	}{
		{name: "water.dat", pfsName: "barren.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "devastation.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "elddar.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "jardelshook.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "maidensgrave.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "mesa.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "moors.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "neighborhood.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "nektulos.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "oceangreenvillage.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "oceanoftears.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "steppes.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "sunderock.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "thalassius.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "theater.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "zhisza.eqg", quadsPerTile: 16},
		// {name: "water.dat", pfsName: "cryptofshade.eqg", quadsPerTile: 16},
		// {name: "bakup_water.dat", pfsName: "buriedsea.eqg", quadsPerTile: 16, isDump: false}, // PASS
		// {name: "bakup_water.dat", pfsName: "jardelshook.eqg", quadsPerTile: 16},              // PASS
		// {name: "bakup_water.dat", pfsName: "maidensgrave.eqg", quadsPerTile: 16},             // PASS
		// {name: "bakup_water.dat", pfsName: "monkeyrock.eqg", quadsPerTile: 16},               // PASS
		// {name: "water.dat", pfsName: "buriedsea.eqg", quadsPerTile: 16},
	}

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

			dat := &DatWtr{}

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
			srcData := data
			dstData := buf.Bytes()
			err = os.WriteFile(fmt.Sprintf("%s/%s.src%s", dirTest, "water", ".dat"), srcData, 0644)
			if err != nil {
				t.Fatalf("failed to write %s: %s", "src water", err.Error())
			}

			err = os.WriteFile(fmt.Sprintf("%s/%s.dst%s", dirTest, "water", ".dat"), dstData, 0644)
			if err != nil {
				t.Fatalf("failed to write %s: %s", "dst water", err.Error())
			}

			dat2 := &DatWtr{}
			err = dat2.Read(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("failed to decode2 %s: %s", tt.name, err.Error())
			}

			err = helper.ByteCompareTest(srcData, dstData)
			if err != nil {
				t.Fatalf("%s byteCompare: %s", tt.name, err)
				return
			}
			fmt.Println(tt.pfsName, tt.name, "PASS")
		})
	}
}
