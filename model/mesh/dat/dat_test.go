package dat

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/tag"
)

func TestDecode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)

	tests := []struct {
		name    string
		pfsName string
		wantErr bool
	}{
		// .dat|?|arthicrex_te.dat|arthicrex.eqg
		// .dat|?|ascent.dat|direwind.eqg
		// .dat|?|bakup_water.dat|buriedsea.eqg
		// .dat|?|bakup_water.dat|jardelshook.eqg
		// .dat|?|bakup_water.dat|maidensgrave.eqg
		// .dat|?|bakup_water.dat|monkeyrock.eqg
		// .dat|?|barrencoast.dat|barren.eqg
		// .dat|?|blacksail.dat|blacksail.eqg
		// .dat|?|brotherisland.dat|brotherisland.eqg
		// .dat|?|buriedsea.dat|buriedsea.eqg
		// .dat|?|commonlands.dat|commonlands.eqg
		// .dat|?|commonlands.dat|oldcommons.eqg
		{name: "cryptofshade.dat", pfsName: "cryptofshade.eqg"},
		// .dat|?|devastation.dat|devastation.eqg
		// .dat|?|dragonscale.dat|dragonscale.eqg
		// .dat|?|empyr.dat|empyr.eqg
		// .dat|?|farstone.dat|arcstone.eqg
		// .dat|?|feerrott.dat|feerrott2.eqg
		// .dat|?|fieldofbone.dat|oldfieldofboneb.eqg
		// .dat|?|floraexclusion.dat|arcstone.eqg
		// .dat|?|floraexclusion.dat|arthicrex.eqg
		// .dat|?|floraexclusion.dat|barren.eqg
		// .dat|?|floraexclusion.dat|blacksail.eqg
		// .dat|?|floraexclusion.dat|brotherisland.eqg
		// .dat|?|floraexclusion.dat|buriedsea.eqg
		// .dat|?|floraexclusion.dat|deadhills.eqg
		// .dat|?|floraexclusion.dat|devastation.eqg
		// .dat|?|floraexclusion.dat|dragonscale.eqg
		// .dat|?|floraexclusion.dat|elddar.eqg
		// .dat|?|floraexclusion.dat|empyr.eqg
		// .dat|?|floraexclusion.dat|frontiermtnsb.eqg
		// .dat|?|floraexclusion.dat|lceanium.eqg
		// .dat|?|floraexclusion.dat|maidensgrave.eqg
		// .dat|?|floraexclusion.dat|mesa.eqg
		// .dat|?|floraexclusion.dat|mistythicket.eqg
		// .dat|?|floraexclusion.dat|neighborhood.eqg
		// .dat|?|floraexclusion.dat|nektulos.eqg
		// .dat|?|floraexclusion.dat|oceangreenhills.eqg
		// .dat|?|floraexclusion.dat|oceangreenvillage.eqg
		// .dat|?|floraexclusion.dat|oceanoftears.eqg
		// .dat|?|floraexclusion.dat|oldbloodfield.eqg
		// .dat|?|floraexclusion.dat|oldkithicor.eqg
		// .dat|?|floraexclusion.dat|overtheretwo.eqg
		// .dat|?|floraexclusion.dat|scorchedwoods.eqg
		// .dat|?|floraexclusion.dat|skyfiretwo.eqg
		// .dat|?|floraexclusion.dat|steppes.eqg
		// .dat|?|floraexclusion.dat|sunderock.eqg
		// .dat|?|floraexclusion.dat|thalassius.eqg
		// .dat|?|floraexclusion.dat|theater.eqg
		// .dat|?|floraexclusion.dat|zhisza.eqg
		// .dat|?|hillsofshade.dat|cryptofshade.eqg
		// .dat|?|innothule.dat|innothuleb.eqg
		// .dat|?|invw.dat|arcstone.eqg
		// .dat|?|invw.dat|arthicrex.eqg
		// .dat|?|invw.dat|blacksail.eqg
		// .dat|?|invw.dat|buriedsea.eqg
		// .dat|?|invw.dat|cryptofshade.eqg
		// .dat|?|invw.dat|deadhills.eqg
		// .dat|?|invw.dat|devastation.eqg
		// .dat|?|invw.dat|direwind.eqg
		// .dat|?|invw.dat|dragonscale.eqg
		// .dat|?|invw.dat|elddar.eqg
		// .dat|?|invw.dat|frontiermtnsb.eqg
		// .dat|?|invw.dat|fungalforest.eqg
		// .dat|?|invw.dat|maidensgrave.eqg
		// .dat|?|invw.dat|mesa.eqg
		// .dat|?|invw.dat|mistythicket.eqg
		// .dat|?|invw.dat|neighborhood.eqg
		// .dat|?|invw.dat|oceangreenhills.eqg
		// .dat|?|invw.dat|oceanoftears.eqg
		// .dat|?|invw.dat|oldfieldofboneb.eqg
		// .dat|?|invw.dat|overtheretwo.eqg
		// .dat|?|invw.dat|scorchedwoods.eqg
		// .dat|?|invw.dat|shadowedmount.eqg
		// .dat|?|invw.dat|shardslanding.eqg
		// .dat|?|invw.dat|skyfiretwo.eqg
		// .dat|?|invw.dat|steamfontmts.eqg
		// .dat|?|invw.dat|steppes.eqg
		// .dat|?|invw.dat|suncrest.eqg
		// .dat|?|invw.dat|sunderock.eqg
		// .dat|?|invw.dat|thalassius.eqg
		// .dat|?|invw.dat|theater.eqg
		// .dat|?|invw.dat|zhisza.eqg
		// .dat|?|lowlands.dat|sunderock.eqg
		// .dat|?|mistythicket.dat|mistythicket.eqg
		// .dat|?|neighborhood.dat|neighborhood.eqg
		// .dat|?|nektuloseditor.dat|nektulos.eqg
		// .dat|?|oceangreenvillage.dat|oceangreenvillage.eqg
		// .dat|?|oldbloodfield.dat|precipiceofwar.eqg
		// .dat|?|oot.dat|oceanoftears.eqg
		// .dat|?|overthere.dat|overtheretwo.eqg
		// .dat|?|planeofmusic.dat|theater.eqg
		// .dat|?|scorchedwoods_terrain.dat|scorchedwoods.eqg
		// .dat|?|skyfiretwo.dat|skyfiretwo.eqg
		// .dat|?|steamfontmts.dat|steamfontmts.eqg
		// .dat|?|suncrest.dat|suncrest.eqg
		// .dat|?|water.dat|arcstone.eqg
		// .dat|?|water.dat|buriedsea.eqg
		// .dat|?|water.dat|cryptofshade.eqg
		// .dat|?|water.dat|deadhills.eqg
		// .dat|?|water.dat|devastation.eqg
		// .dat|?|water.dat|elddar.eqg
		// .dat|?|water.dat|jardelshook.eqg
		// .dat|?|water.dat|maidensgrave.eqg
		// .dat|?|water.dat|mesa.eqg
		// .dat|?|water.dat|moors.eqg
		// .dat|?|water.dat|neighborhood.eqg
		// .dat|?|water.dat|nektulos.eqg
		// .dat|?|water.dat|oceangreenvillage.eqg
		// .dat|?|water.dat|oceanoftears.eqg
		// .dat|?|water.dat|steppes.eqg
		// .dat|?|water.dat|sunderock.eqg
		// .dat|?|water.dat|thalassius.eqg
		// .dat|?|water.dat|theater.eqg
		// .dat|?|water.dat|zhisza.eqg
		// .dat|?|zihssa.dat|zhisza.eqg
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.pfsName))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".zon" {
					continue
				}
				zone := common.NewZone("")

				err = Decode(zone, bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
				if err != nil {
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
				}

			}
		})
	}
}

func TestEncode(t *testing.T) {
	eqPath := os.Getenv("EQ_PATH")
	if eqPath == "" {
		t.Skip("EQ_PATH not set")
	}
	dirTest := common.DirTest(t)

	tests := []struct {
		name    string
		pfsName string
		wantErr bool
	}{
		{name: "cryptofshade.dat", pfsName: "cryptofshade.eqg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfs, err := pfs.NewFile(fmt.Sprintf("%s/%s", eqPath, tt.pfsName))
			if err != nil {
				t.Fatalf("failed to open eqg %s: %s", tt.name, err.Error())
			}
			for _, file := range pfs.Files() {
				if filepath.Ext(file.Name()) != ".zon" {
					continue
				}
				zone := common.NewZone("")

				err = Decode(zone, bytes.NewReader(file.Data()))
				os.WriteFile(fmt.Sprintf("%s/%s", dirTest, file.Name()), file.Data(), 0644)
				tag.Write(fmt.Sprintf("%s/%s.tags", dirTest, file.Name()))
				if err != nil {
					t.Fatalf("failed to decode %s: %s", tt.name, err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = Encode(zone, buf)
				if err != nil {
					t.Fatalf("failed to encode %s: %s", tt.name, err.Error())
				}

				//srcData := file.Data()
				//dstData := buf.Bytes()
				/*for i := 0; i < len(srcData); i++ {
					if len(dstData) <= i {
						min := 0
						max := len(srcData)
						fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
						max = len(dstData)
						fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))

						t.Fatalf("%s src eof at offset %d (dst is too large by %d bytes)", tt.name, i, len(dstData)-len(srcData))
					}
					if len(dstData) <= i {
						t.Fatalf("%s dst eof at offset %d (dst is too small by %d bytes)", tt.name, i, len(srcData)-len(dstData))
					}
					if srcData[i] == dstData[i] {
						continue
					}

					fmt.Printf("%s mismatch at offset %d (src: 0x%x vs dst: 0x%x aka %d)\n", tt.name, i, srcData[i], dstData[i], dstData[i])
					max := i + 16
					if max > len(srcData) {
						max = len(srcData)
					}

					min := i - 16
					if min < 0 {
						min = 0
					}
					fmt.Printf("src (%d:%d):\n%s\n", min, max, hex.Dump(srcData[min:max]))
					if max > len(dstData) {
						max = len(dstData)
					}

					fmt.Printf("dst (%d:%d):\n%s\n", min, max, hex.Dump(dstData[min:max]))
					//os.WriteFile(fmt.Sprintf("%s/_src_%s", dirTest, file.Name()), file.Data(), 0644)
					//os.WriteFile(fmt.Sprintf("%s/_dst_%s", dirTest, file.Name()), buf.Bytes(), 0644)
					t.Fatalf("%s encode: data mismatch", tt.name)
				}*/
			}
		})
	}
}
