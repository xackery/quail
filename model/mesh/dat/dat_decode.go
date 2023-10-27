package dat

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/tag"
)

// Decode decodes a v4 zone dat file
// https://github.com/EQEmu/zone-utilities/blob/master/src/common/eqg_v4_loader.cpp#L115
func Decode(zone *common.Zone, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	tag.New()

	zone.Header.Version = 4
	dst := zone.V4Dat

	dst.Unk1 = dec.Uint32()
	dst.Unk2 = dec.Uint32()
	dst.Unk3 = dec.Uint32()
	dst.BaseTileTexture = dec.StringZero()
	tileCount := dec.Uint32()

	//zoneMinX := float32(float32(zone.V4Info.MinLat*zone.V4Info.QuadsPerTile) * zone.V4Info.UnitsPerVert)
	//zoneMinY := float32(float32(zone.V4Info.MinLng*zone.V4Info.QuadsPerTile) * zone.V4Info.UnitsPerVert)
	quadCount := zone.V4Info.QuadsPerTile * zone.V4Info.QuadsPerTile
	vertCount := (zone.V4Info.QuadsPerTile + 1) * (zone.V4Info.QuadsPerTile + 1)
	for i := 0; i < int(tileCount); i++ {
		tile := common.V4DatTile{}
		tile.Lng = dec.Int32()
		tile.Lat = dec.Int32()
		tile.Unk = dec.Uint32()
		//tileStartX := float32(zoneMinX + (float32(tileLat) - 100000 - float32(zone.V4Info.MinLat)*float32(zone.V4Info.UnitsPerVert)*float32(zone.V4Info.QuadsPerTile)))
		//tileStartY := float32(zoneMinY + (float32(tileLng) - 100000 - float32(zone.V4Info.MinLng)*float32(zone.V4Info.UnitsPerVert)*float32(zone.V4Info.QuadsPerTile)))

		//isFloatsAllSame := true

		for j := 0; j < vertCount; j++ {
			tile.Colors = append(tile.Colors, dec.Uint32())
		}

		for j := 0; j < vertCount; j++ {
			tile.Colors2 = append(tile.Colors2, dec.Uint32())
		}

		flags := []uint8{}
		for j := 0; j < quadCount; j++ {
			flag := dec.Uint8()
			flags = append(flags, flag)
			if flag&0x01 == 0x01 {
				//isFloatsAllSame = false
			}
		}
		/*
			isFlat := isFloatsAllSame

			baseWaterLevel := dec.Float32()

			unk1 := dec.Int32()
			unk2 := dec.Int32()

			if unk1 > 0 {

			}
		*/
		dst.Tiles = append(dst.Tiles, tile)
	}

	zone.V4Dat = dst
	return nil
}
