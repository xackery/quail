package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

type DatZon struct {
	MetaFileName          string
	Version               uint32
	Flags                 uint32
	FallbackDetailRepeat  uint32
	Unusued               uint32
	FallbackDetailMapName string
	QuadsPerTile          int
	Tiles                 []*DatZonTile
}

func (dat *DatZon) Identity() string {
	return "dat"
}

type DatZonTile struct {
	Lng               int32
	Lat               int32
	Unk1              uint32
	Floats            []float32
	Colors            []uint32
	Colors2           []uint32
	Flags             []uint8
	BaseWaterLevel    float32
	Unk2              int32
	Unk3              int8
	Unk3Quad          [4]float32
	Unk3Float         float32
	LayerBaseMaterial string
	Layers            []*DatZonLayer
	SinglePlacables   []*DatZonSinglePlacable
	Areas             []*DatZonArea
	LightEffects      []*DatZonLightEffect
	TogRefs           []*DatZonTogRef
}

type DatZonLayer struct {
	Material       string
	DetailMaskDim  uint32
	DetailMaskDims []uint8
}

type DatZonSinglePlacable struct {
	ModelName    string
	InstanceName string
	Longitude    int32
	Latitude     int32

	Position [3]float32
	Rotation [3]float32
	Scale    [3]float32
	Flags    uint8
	Unk1     uint32
}

type DatZonArea struct {
	UnkStr1   string
	Type      int32
	UnkStr2   string
	Longitude uint32
	Latitude  uint32
	Position  [3]float32
	Rotation  [3]float32
	Scale     [3]float32
	Size      [3]float32
}

type DatZonLightEffect struct {
	UnkStr1   string
	UnkStr2   string
	Unk3      uint8
	Longitude uint32
	Latitude  uint32
	Position  [3]float32
	Rotation  [3]float32
	Scale     [3]float32
	Size      float32
}

type DatZonTogRef struct {
	Name      string
	Longitude uint32
	Latitude  uint32
	Position  [3]float32
	Rotation  [3]float32
	Scale     [3]float32
	Adjust    float32
}

// Decode reads a v4 zone dat file
// https://github.com/EQEmu/zone-utilities/blob/master/src/common/eqg_v4_loader.cpp#L115
func (e *DatZon) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	e.Version = dec.Uint32()

	e.Flags = dec.Uint32()
	e.FallbackDetailRepeat = dec.Uint32()
	e.FallbackDetailMapName = dec.StringZero()
	tileCount := dec.Uint32()

	//zoneMinX := float32(float32(dat.MinLat*dat.QuadsPerTile) * dat.UnitsPerVert)
	//zoneMinY := float32(float32(dat.MinLng*dat.QuadsPerTile) * dat.UnitsPerVert)
	quadCount := e.QuadsPerTile * e.QuadsPerTile
	vertCount := (e.QuadsPerTile + 1) * (e.QuadsPerTile + 1)
	if tileCount > 9999 {
		return fmt.Errorf("tile count %d is too high", tileCount)
	}
	for i := 0; i < int(tileCount); i++ {
		tile := &DatZonTile{}
		tile.Lng = dec.Int32()
		tile.Lat = dec.Int32()
		tile.Unk1 = dec.Uint32()
		//tileStartX := float32(zoneMinX + (float32(tileLat) - 100000 - float32(dat.MinLat)*float32(dat.UnitsPerVert)*float32(dat.QuadsPerTile)))
		//tileStartY := float32(zoneMinY + (float32(tileLng) - 100000 - float32(dat.MinLng)*float32(dat.UnitsPerVert)*float32(dat.QuadsPerTile)))

		//isFloatsAllSame := true

		for j := 0; j < vertCount; j++ {
			tile.Floats = append(tile.Floats, dec.Float32())
		}

		for j := 0; j < vertCount; j++ {
			tile.Colors = append(tile.Colors, dec.Uint32())
		}

		for j := 0; j < vertCount; j++ {
			tile.Colors2 = append(tile.Colors2, dec.Uint32())
		}

		for j := 0; j < quadCount; j++ {
			flag := dec.Uint8()
			//if flag&0x01 == 0x01 {
			//isFloatsAllSame = false
			//}
			tile.Flags = append(tile.Flags, flag)
		}
		//isFlat := isFloatsAllSame

		tile.BaseWaterLevel = dec.Float32()
		tile.Unk2 = dec.Int32()

		if tile.Unk2 > 0 {
			tile.Unk3 = dec.Int8()
			if tile.Unk3 > 0 {
				tile.Unk3Quad[0] = dec.Float32()
				tile.Unk3Quad[1] = dec.Float32()
				tile.Unk3Quad[2] = dec.Float32()
				tile.Unk3Quad[3] = dec.Float32()
			}
			tile.Unk3Float = dec.Float32()
		}

		layerCount := dec.Uint32()
		if layerCount > 0 {
			tile.LayerBaseMaterial = dec.StringZero()
		}
		if layerCount > 9999 {
			return fmt.Errorf("layer count %d is too high", layerCount)
		}
		for j := 1; j < int(layerCount); j++ {
			layer := &DatZonLayer{}
			layer.Material = dec.StringZero()
			layer.DetailMaskDim = dec.Uint32()

			szM := layer.DetailMaskDim * layer.DetailMaskDim
			for k := 0; k < int(szM); k++ {
				layer.DetailMaskDims = append(layer.DetailMaskDims, dec.Uint8())
			}

			tile.Layers = append(tile.Layers, layer)
		}

		singlePlacableCount := dec.Uint32()
		if singlePlacableCount > 9999 {
			return fmt.Errorf("single placable count %d is too high", singlePlacableCount)
		}
		for j := 0; j < int(singlePlacableCount); j++ {
			singlePlacable := &DatZonSinglePlacable{}
			singlePlacable.ModelName = dec.StringZero()
			singlePlacable.InstanceName = dec.StringZero()
			singlePlacable.Longitude = dec.Int32()
			singlePlacable.Latitude = dec.Int32()
			singlePlacable.Position[0] = dec.Float32()
			singlePlacable.Position[1] = dec.Float32()
			singlePlacable.Position[2] = dec.Float32()
			singlePlacable.Rotation[0] = dec.Float32()
			singlePlacable.Rotation[1] = dec.Float32()
			singlePlacable.Rotation[2] = dec.Float32()
			singlePlacable.Scale[0] = dec.Float32()
			singlePlacable.Scale[1] = dec.Float32()
			singlePlacable.Scale[2] = dec.Float32()
			singlePlacable.Flags = dec.Uint8()

			if e.Flags&0x02 == 2 {
				singlePlacable.Unk1 = dec.Uint32()
			}

			tile.SinglePlacables = append(tile.SinglePlacables, singlePlacable)
		}

		areasCount := dec.Uint32()
		for j := 0; j < int(areasCount); j++ {
			area := &DatZonArea{}
			area.UnkStr1 = dec.StringZero()
			area.Type = dec.Int32()
			area.UnkStr2 = dec.StringZero()
			area.Longitude = dec.Uint32()
			area.Latitude = dec.Uint32()
			area.Position[0] = dec.Float32()
			area.Position[1] = dec.Float32()
			area.Position[2] = dec.Float32()
			area.Rotation[0] = dec.Float32()
			area.Rotation[1] = dec.Float32()
			area.Rotation[2] = dec.Float32()
			area.Scale[0] = dec.Float32()
			area.Scale[1] = dec.Float32()
			area.Scale[2] = dec.Float32()
			area.Size[0] = dec.Float32()
			area.Size[1] = dec.Float32()
			area.Size[2] = dec.Float32()
			tile.Areas = append(tile.Areas, area)
		}

		lightEffectsCount := dec.Uint32()
		for j := 0; j < int(lightEffectsCount); j++ {
			lightEffect := &DatZonLightEffect{}
			lightEffect.UnkStr1 = dec.StringZero()
			lightEffect.UnkStr2 = dec.StringZero()
			lightEffect.Unk3 = dec.Uint8()
			lightEffect.Longitude = dec.Uint32()
			lightEffect.Latitude = dec.Uint32()
			lightEffect.Position[0] = dec.Float32()
			lightEffect.Position[1] = dec.Float32()
			lightEffect.Position[2] = dec.Float32()
			lightEffect.Rotation[0] = dec.Float32()
			lightEffect.Rotation[1] = dec.Float32()
			lightEffect.Rotation[2] = dec.Float32()
			lightEffect.Scale[0] = dec.Float32()
			lightEffect.Scale[1] = dec.Float32()
			lightEffect.Scale[2] = dec.Float32()
			lightEffect.Size = dec.Float32()
			tile.LightEffects = append(tile.LightEffects, lightEffect)
		}

		togRefsCount := dec.Uint32()
		for j := 0; j < int(togRefsCount); j++ {
			togRef := &DatZonTogRef{}
			togRef.Name = dec.StringZero()
			togRef.Longitude = dec.Uint32()
			togRef.Latitude = dec.Uint32()
			togRef.Position[0] = dec.Float32()
			togRef.Position[1] = dec.Float32()
			togRef.Position[2] = dec.Float32()
			togRef.Rotation[0] = dec.Float32()
			togRef.Rotation[1] = dec.Float32()
			togRef.Rotation[2] = dec.Float32()
			togRef.Scale[0] = dec.Float32()
			togRef.Scale[1] = dec.Float32()
			togRef.Scale[2] = dec.Float32()
			togRef.Adjust = dec.Float32()
			tile.TogRefs = append(tile.TogRefs, togRef)
		}

		e.Tiles = append(e.Tiles, tile)
	}

	pos := dec.Pos()
	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if pos != endPos {
		if pos < endPos {
			return fmt.Errorf("%d bytes remaining (%d total)", endPos-pos, endPos)
		}

		return fmt.Errorf("read past end of file")
	}
	return nil
}

// SetName sets the name of the file
func (e *DatZon) SetFileName(name string) {
	e.MetaFileName = name
}

func (e *DatZon) FileName() string {
	return e.MetaFileName
}
