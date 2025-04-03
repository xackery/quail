package raw

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestDatIwWrite(t *testing.T) {
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
		// .dat|?|invw.dat|arcstone.eqg
		{name: "invw.dat", pfsName: "arcstone.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|arthicrex.eqg
		//{name: "invw.dat", pfsName: "arthicrex.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|blacksail.eqg
		//{name: "invw.dat", pfsName: "blacksail.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|buriedsea.eqg
		//{name: "invw.dat", pfsName: "buriedsea.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|cryptofshade.eqg
		//{name: "invw.dat", pfsName: "cryptofshade.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|deadhills.eqg
		////{name: "invw.dat", pfsName: "deadhills.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|invw.dat|devastation.eqg
		//{name: "invw.dat", pfsName: "devastation.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|direwind.eqg
		//{name: "invw.dat", pfsName: "direwind.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|dragonscale.eqg
		//{name: "invw.dat", pfsName: "dragonscale.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|elddar.eqg
		//{name: "invw.dat", pfsName: "elddar.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|frontiermtnsb.eqg
		//{name: "invw.dat", pfsName: "frontiermtnsb.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|fungalforest.eqg
		//{name: "invw.dat", pfsName: "fungalforest.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|maidensgrave.eqg
		//{name: "invw.dat", pfsName: "maidensgrave.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|mesa.eqg
		// {name: "invw.dat", pfsName: "mesa.eqg", quadsPerTile: 16}, // too high
		// .dat|?|invw.dat|mistythicket.eqg
		//{name: "invw.dat", pfsName: "mistythicket.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|neighborhood.eqg
		//{name: "invw.dat", pfsName: "neighborhood.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|oceangreenhills.eqg
		//{name: "invw.dat", pfsName: "oceangreenhills.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|oceanoftears.eqg
		//{name: "invw.dat", pfsName: "oceanoftears.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|oldfieldofboneb.eqg
		//{name: "invw.dat", pfsName: "oldfieldofboneb.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|overtheretwo.eqg
		//{name: "invw.dat", pfsName: "overtheretwo.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|scorchedwoods.eqg
		//{name: "invw.dat", pfsName: "scorchedwoods.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|shadowedmount.eqg
		//{name: "invw.dat", pfsName: "shadowedmount.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|shardslanding.eqg
		//{name: "invw.dat", pfsName: "shardslanding.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|skyfiretwo.eqg
		//{name: "invw.dat", pfsName: "skyfiretwo.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|steamfontmts.eqg
		//{name: "invw.dat", pfsName: "steamfontmts.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|steppes.eqg
		//{name: "invw.dat", pfsName: "steppes.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|suncrest.eqg
		//{name: "invw.dat", pfsName: "suncrest.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|sunderock.eqg
		//{name: "invw.dat", pfsName: "sunderock.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|thalassius.eqg
		//{name: "invw.dat", pfsName: "thalassius.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|theater.eqg
		//{name: "invw.dat", pfsName: "theater.eqg", quadsPerTile: 16},
		// .dat|?|invw.dat|zhisza.eqg
		//{name: "invw.dat", pfsName: "zhisza.eqg", quadsPerTile: 16},
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
