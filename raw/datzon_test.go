package raw

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
)

func TestDatWrite(t *testing.T) {
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
		//{name: "hillsofshade.dat", pfsName: "cryptofshade.eqg", quadsPerTile: 16, isDump: false}, // PASS
		// .dat|?|ascent.dat|direwind.eqg
		//{name: "ascent.dat", pfsName: "direwind.eqg", quadsPerTile: 16, isDump: false}, // PASS
		// .dat|?|barrencoast.dat|barren.eqg
		//{name: "barrencoast.dat", pfsName: "barren.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|blacksail.dat|blacksail.eqg
		//{name: "blacksail.dat", pfsName: "blacksail.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|brotherisland.dat|brotherisland.eqg
		//{name: "brotherisland.dat", pfsName: "brotherisland.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|buriedsea.dat|buriedsea.eqg
		//{name: "buriedsea.dat", pfsName: "buriedsea.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|commonlands.dat|commonlands.eqg
		//{name: "commonlands.dat", pfsName: "commonlands.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|commonlands.dat|oldcommons.eqg
		//{name: "commonlands.dat", pfsName: "oldcommons.eqg", quadsPerTile: 16}, // TODO: FAIL
		// .dat|?|devastation.dat|devastation.eqg
		// {name: "devastation.dat", pfsName: "devastation.eqg", quadsPerTile: 16}, // TODO: fix
		// .dat|?|dragonscale.dat|dragonscale.eqg
		//{name: "dragonscale.dat", pfsName: "dragonscale.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|empyr.dat|empyr.eqg
		//{name: "empyr.dat", pfsName: "empyr.eqg", quadsPerTile: 16}, // TODO: fix
		// .dat|?|farstone.dat|arcstone.eqg
		//{name: "farstone.dat", pfsName: "arcstone.eqg", quadsPerTile: 16}, // TODO: fix
		// .dat|?|feerrott.dat|feerrott2.eqg
		//{name: "feerrott.dat", pfsName: "feerrott2.eqg", quadsPerTile: 16},
		// .dat|?|fieldofbone.dat|oldfieldofboneb.eqg
		//{name: "fieldofbone.dat", pfsName: "oldfieldofboneb.eqg", quadsPerTile: 16},
		// .dat|?|floraexclusion.dat|arcstone.eqg
		//{name: "floraexclusion.dat", pfsName: "arcstone.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|arthicrex.eqg
		//{name: "floraexclusion.dat", pfsName: "arthicrex.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|barren.eqg
		//{name: "floraexclusion.dat", pfsName: "barren.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|blacksail.eqg
		//{name: "floraexclusion.dat", pfsName: "blacksail.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|brotherisland.eqg
		//{name: "floraexclusion.dat", pfsName: "brotherisland.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|buriedsea.eqg
		//{name: "floraexclusion.dat", pfsName: "buriedsea.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|deadhills.eqg
		//{name: "floraexclusion.dat", pfsName: "deadhills.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|devastation.eqg
		//{name: "floraexclusion.dat", pfsName: "devastation.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|dragonscale.eqg
		//{name: "floraexclusion.dat", pfsName: "dragonscale.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|elddar.eqg
		//{name: "floraexclusion.dat", pfsName: "elddar.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|empyr.eqg
		//{name: "floraexclusion.dat", pfsName: "empyr.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|frontiermtnsb.eqg
		//{name: "floraexclusion.dat", pfsName: "frontiermtnsb.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|lceanium.eqg
		//{name: "floraexclusion.dat", pfsName: "lceanium.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|maidensgrave.eqg
		//{name: "floraexclusion.dat", pfsName: "maidensgrave.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|mesa.eqg
		//{name: "floraexclusion.dat", pfsName: "mesa.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|mistythicket.eqg
		//{name: "floraexclusion.dat", pfsName: "mistythicket.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|neighborhood.eqg
		//{name: "floraexclusion.dat", pfsName: "neighborhood.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|nektulos.eqg
		//{name: "floraexclusion.dat", pfsName: "nektulos.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|oceangreenhills.eqg
		//{name: "floraexclusion.dat", pfsName: "oceangreenhills.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|oceangreenvillage.eqg
		//{name: "floraexclusion.dat", pfsName: "oceangreenvillage.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|oceanoftears.eqg
		//{name: "floraexclusion.dat", pfsName: "oceanoftears.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|oldbloodfield.eqg
		//{name: "floraexclusion.dat", pfsName: "oldbloodfield.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|oldkithicor.eqg
		//{name: "floraexclusion.dat", pfsName: "oldkithicor.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|overtheretwo.eqg
		//{name: "floraexclusion.dat", pfsName: "overtheretwo.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|scorchedwoods.eqg
		//{name: "floraexclusion.dat", pfsName: "scorchedwoods.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|skyfiretwo.eqg
		//{name: "floraexclusion.dat", pfsName: "skyfiretwo.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|steppes.eqg
		//{name: "floraexclusion.dat", pfsName: "steppes.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|sunderock.eqg
		//{name: "floraexclusion.dat", pfsName: "sunderock.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|thalassius.eqg
		//{name: "floraexclusion.dat", pfsName: "thalassius.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|theater.eqg
		//{name: "floraexclusion.dat", pfsName: "theater.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|floraexclusion.dat|zhisza.eqg
		//{name: "floraexclusion.dat", pfsName: "zhisza.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|hillsofshade.dat|cryptofshade.eqg
		//{name: "hillsofshade.dat", pfsName: "cryptofshade.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|innothule.dat|innothuleb.eqg
		//{name: "innothule.dat", pfsName: "innothuleb.eqg", quadsPerTile: 16}, // PASS
		// .dat|?|invw.dat|arcstone.eqg
		//{name: "invw.dat", pfsName: "arcstone.eqg", quadsPerTile: 16},
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
		// .dat|?|lowlands.dat|sunderock.eqg
		//{name: "lowlands.dat", pfsName: "sunderock.eqg", quadsPerTile: 16},
		// .dat|?|mistythicket.dat|mistythicket.eqg
		//{name: "mistythicket.dat", pfsName: "mistythicket.eqg", quadsPerTile: 16},
		// .dat|?|neighborhood.dat|neighborhood.eqg
		//{name: "neighborhood.dat", pfsName: "neighborhood.eqg", quadsPerTile: 16},
		// .dat|?|nektuloseditor.dat|nektulos.eqg
		//{name: "nektuloseditor.dat", pfsName: "nektulos.eqg", quadsPerTile: 16},
		// .dat|?|oceangreenvillage.dat|oceangreenvillage.eqg
		//{name: "oceangreenvillage.dat", pfsName: "oceangreenvillage.eqg", quadsPerTile: 16},
		// .dat|?|oldbloodfield.dat|precipiceofwar.eqg
		//{name: "oldbloodfield.dat", pfsName: "precipiceofwar.eqg", quadsPerTile: 16},
		// .dat|?|oot.dat|oceanoftears.eqg
		//{name: "oot.dat", pfsName: "oceanoftears.eqg", quadsPerTile: 16},
		// .dat|?|overthere.dat|overtheretwo.eqg
		//{name: "overthere.dat", pfsName: "overtheretwo.eqg", quadsPerTile: 16},
		// .dat|?|planeofmusic.dat|theater.eqg
		//{name: "planeofmusic.dat", pfsName: "theater.eqg", quadsPerTile: 16},
		// .dat|?|scorchedwoods_terrain.dat|scorchedwoods.eqg
		//{name: "scorchedwoods_terrain.dat", pfsName: "scorchedwoods.eqg", quadsPerTile: 16},
		// .dat|?|skyfiretwo.dat|skyfiretwo.eqg
		//{name: "skyfiretwo.dat", pfsName: "skyfiretwo.eqg", quadsPerTile: 16},
		// .dat|?|steamfontmts.dat|steamfontmts.eqg
		//{name: "steamfontmts.dat", pfsName: "steamfontmts.eqg", quadsPerTile: 16},
		// .dat|?|suncrest.dat|suncrest.eqg
		//{name: "suncrest.dat", pfsName: "suncrest.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|arcstone.eqg
		//{name: "water.dat", pfsName: "arcstone.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|buriedsea.eqg
		//{name: "water.dat", pfsName: "buriedsea.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|cryptofshade.eqg
		//{name: "water.dat", pfsName: "cryptofshade.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|deadhills.eqg
		//{name: "water.dat", pfsName: "deadhills.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|devastation.eqg
		//{name: "water.dat", pfsName: "devastation.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|elddar.eqg
		//{name: "water.dat", pfsName: "elddar.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|jardelshook.eqg
		//{name: "water.dat", pfsName: "jardelshook.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|maidensgrave.eqg
		//{name: "water.dat", pfsName: "maidensgrave.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|mesa.eqg
		//{name: "water.dat", pfsName: "mesa.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|moors.eqg
		//{name: "water.dat", pfsName: "moors.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|neighborhood.eqg
		//{name: "water.dat", pfsName: "neighborhood.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|nektulos.eqg
		//{name: "water.dat", pfsName: "nektulos.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|oceangreenvillage.eqg
		//{name: "water.dat", pfsName: "oceangreenvillage.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|oceanoftears.eqg
		//{name: "water.dat", pfsName: "oceanoftears.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|steppes.eqg
		//{name: "water.dat", pfsName: "steppes.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|sunderock.eqg
		//{name: "water.dat", pfsName: "sunderock.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|thalassius.eqg
		//{name: "water.dat", pfsName: "thalassius.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|theater.eqg
		//{name: "water.dat", pfsName: "theater.eqg", quadsPerTile: 16},
		// .dat|?|water.dat|zhisza.eqg
		//{name: "water.dat", pfsName: "zhisza.eqg", quadsPerTile: 16},
		// .dat|?|zihssa.dat|zhisza.eqg
		//{name: "zihssa.dat", pfsName: "zhisza.eqg", quadsPerTile: 16},
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

			dat := &DatZon{
				QuadsPerTile: tt.quadsPerTile,
			}

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

			dat2 := &DatZon{
				QuadsPerTile: tt.quadsPerTile,
			}
			err = dat2.Read(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("failed to decode2 %s: %s", tt.name, err.Error())
			}

			if len(dat.Tiles) != len(dat2.Tiles) {
				t.Fatalf("tile count mismatch %d != %d", len(dat.Tiles), len(dat2.Tiles))
			}

			for i := range dat.Tiles {
				tile := dat.Tiles[i]
				tile2 := dat2.Tiles[i]
				if tile.Lng != tile2.Lng {
					t.Fatalf("tile lng mismatch %d != %d", tile.Lng, tile2.Lng)
				}
				if tile.Lat != tile2.Lat {
					t.Fatalf("tile lat mismatch %d != %d", tile.Lat, tile2.Lat)
				}
				if tile.Unk1 != tile2.Unk1 {
					t.Fatalf("tile unk1 mismatch %d != %d", tile.Unk1, tile2.Unk1)
				}
				if len(tile.Colors) != len(tile2.Colors) {
					t.Fatalf("tile color count mismatch %d != %d", len(tile.Colors), len(tile2.Colors))
				}
				if len(tile.Colors2) != len(tile2.Colors2) {
					t.Fatalf("tile color2 count mismatch %d != %d", len(tile.Colors2), len(tile2.Colors2))
				}
				if len(tile.Flags) != len(tile2.Flags) {
					t.Fatalf("tile flag count mismatch %d != %d", len(tile.Flags), len(tile2.Flags))
				}
			}

			err = helper.ByteCompareTest(srcData, dstData)
			if err != nil {
				t.Fatalf("%s byteCompare: %s", tt.name, err)
			}
			fmt.Println(tt.pfsName, tt.name, "PASS")
		})
	}
}
