package raw

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/tag"
)

type Dat struct {
	MetaFileName    string     `yaml:"file_name"`
	Version         uint32     `yaml:"version"`
	Unk1            uint32     `yaml:"unk1"`
	Unk2            uint32     `yaml:"unk2"`
	Unk3            uint32     `yaml:"unk3"`
	BaseTileTexture string     `yaml:"base_tile_texture"`
	QuadsPerTile    int        `yaml:"quads_per_tile"`
	Tiles           []*DatTile `yaml:"tiles"`
}

type DatTile struct {
	Lng     int32    `yaml:"lng,omitempty"`
	Lat     int32    `yaml:"lat,omitempty"`
	Unk     uint32   `yaml:"unk,omitempty"`
	Colors  []uint32 `yaml:"colors,omitempty"`
	Colors2 []uint32 `yaml:"colors2,omitempty"`
}

// Decode reads a v4 zone dat file
// https://github.com/EQEmu/zone-utilities/blob/master/src/common/eqg_v4_loader.cpp#L115
func (dat *Dat) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	tag.New()

	dat.Version = 4

	dat.Unk1 = dec.Uint32()
	dat.Unk2 = dec.Uint32()
	dat.Unk3 = dec.Uint32()
	dat.BaseTileTexture = dec.StringZero()
	tileCount := dec.Uint32()

	//zoneMinX := float32(float32(dat.MinLat*dat.QuadsPerTile) * dat.UnitsPerVert)
	//zoneMinY := float32(float32(dat.MinLng*dat.QuadsPerTile) * dat.UnitsPerVert)
	quadCount := dat.QuadsPerTile * dat.QuadsPerTile
	vertCount := (dat.QuadsPerTile + 1) * (dat.QuadsPerTile + 1)
	for i := 0; i < int(tileCount); i++ {
		tile := &DatTile{}
		tile.Lng = dec.Int32()
		tile.Lat = dec.Int32()
		tile.Unk = dec.Uint32()
		//tileStartX := float32(zoneMinX + (float32(tileLat) - 100000 - float32(dat.MinLat)*float32(dat.UnitsPerVert)*float32(dat.QuadsPerTile)))
		//tileStartY := float32(zoneMinY + (float32(tileLng) - 100000 - float32(dat.MinLng)*float32(dat.UnitsPerVert)*float32(dat.QuadsPerTile)))

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
		dat.Tiles = append(dat.Tiles, tile)
	}

	return nil
}

// SetName sets the name of the file
func (dat *Dat) SetFileName(name string) {
	dat.MetaFileName = name
}

func (dat *Dat) FileName() string {
	return dat.MetaFileName
}
