package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Encode encodes a v4 zone dat file
// https://github.com/EQEmu/zone-utilities/blob/master/src/common/eqg_v4_loader.cpp#L115
func (e *DatZon) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Uint32(e.Flags)
	enc.Uint32(e.FallbackDetailRepeat)
	enc.Uint32(e.Unusued)
	enc.StringZero(e.FallbackDetailMapName)
	enc.Uint32(uint32(len(e.Tiles)))
	for _, tile := range e.Tiles {
		enc.Int32(tile.Lng)
		enc.Int32(tile.Lat)
		enc.Uint32(tile.Unk1)
		for _, f := range tile.Floats {
			enc.Float32(f)
		}
		for _, color := range tile.Colors {
			enc.Uint32(color)
		}
		for _, color := range tile.Colors2 {
			enc.Uint32(color)
		}
		for _, flag := range tile.Flags {
			enc.Uint8(flag)
		}
		enc.Float32(tile.BaseWaterLevel)
		enc.Int32(tile.Unk2)
		if tile.Unk2 > 0 {
			enc.Int8(tile.Unk3)
			if tile.Unk3 > 0 {
				enc.Float32(tile.Unk3Quad[0])
				enc.Float32(tile.Unk3Quad[1])
				enc.Float32(tile.Unk3Quad[2])
				enc.Float32(tile.Unk3Quad[3])
			}
			enc.Float32(tile.Unk3Float)
		}
		enc.Uint32(uint32(len(tile.Layers)))
		if len(tile.Layers) > 0 {
			enc.StringZero(tile.LayerBaseMaterial)
		}

		if len(tile.Layers) > 1 {
			for _, layer := range tile.Layers {
				enc.StringZero(layer.Material)
				enc.Uint32(layer.DetailMaskDim)
				szM := layer.DetailMaskDim * layer.DetailMaskDim
				if szM != uint32(len(layer.DetailMaskDims)) {
					return fmt.Errorf("detail mask dim times itself %d does not match detail mask dims %d", szM, len(layer.DetailMaskDims))
				}
				for _, dims := range layer.DetailMaskDims {
					enc.Uint8(dims)
				}
			}
		}

		enc.Uint32(uint32(len(tile.SinglePlacables)))
		for _, singlePlacable := range tile.SinglePlacables {
			enc.StringZero(singlePlacable.ModelName)
			enc.StringZero(singlePlacable.InstanceName)
			enc.Int32(singlePlacable.Longitude)
			enc.Int32(singlePlacable.Latitude)
			enc.Float32(singlePlacable.Position[0])
			enc.Float32(singlePlacable.Position[1])
			enc.Float32(singlePlacable.Position[2])
			enc.Float32(singlePlacable.Rotation[0])
			enc.Float32(singlePlacable.Rotation[1])
			enc.Float32(singlePlacable.Rotation[2])
			enc.Float32(singlePlacable.Scale[0])
			enc.Float32(singlePlacable.Scale[1])
			enc.Float32(singlePlacable.Scale[2])
			enc.Uint8(singlePlacable.Flags)
			if e.Flags&0x02 == 2 {
				enc.Uint32(singlePlacable.Unk1)
			}
		}

		enc.Uint32(uint32(len(tile.Areas)))
		for _, area := range tile.Areas {
			enc.StringZero(area.UnkStr1)
			enc.Int32(area.Type)
			enc.StringZero(area.UnkStr2)
			enc.Uint32(area.Longitude)
			enc.Uint32(area.Latitude)
			enc.Float32(area.Position[0])
			enc.Float32(area.Position[1])
			enc.Float32(area.Position[2])
			enc.Float32(area.Rotation[0])
			enc.Float32(area.Rotation[1])
			enc.Float32(area.Rotation[2])
			enc.Float32(area.Scale[0])
			enc.Float32(area.Scale[1])
			enc.Float32(area.Scale[2])
			enc.Float32(area.Size[0])
			enc.Float32(area.Size[1])
			enc.Float32(area.Size[2])
		}

		enc.Uint32(uint32(len(tile.LightEffects)))
		for _, lightEffect := range tile.LightEffects {
			enc.StringZero(lightEffect.UnkStr1)
			enc.StringZero(lightEffect.UnkStr2)
			enc.Uint8(lightEffect.Unk3)
			enc.Uint32(lightEffect.Longitude)
			enc.Uint32(lightEffect.Latitude)
			enc.Float32(lightEffect.Position[0])
			enc.Float32(lightEffect.Position[1])
			enc.Float32(lightEffect.Position[2])
			enc.Float32(lightEffect.Rotation[0])
			enc.Float32(lightEffect.Rotation[1])
			enc.Float32(lightEffect.Rotation[2])
			enc.Float32(lightEffect.Scale[0])
			enc.Float32(lightEffect.Scale[1])
			enc.Float32(lightEffect.Scale[2])
			enc.Float32(lightEffect.Size)
		}

		enc.Uint32(uint32(len(tile.TogRefs)))
		for _, datTogRef := range tile.TogRefs {
			enc.StringZero(datTogRef.Name)
			enc.Uint32(datTogRef.Longitude)
			enc.Uint32(datTogRef.Latitude)
			enc.Float32(datTogRef.Position[0])
			enc.Float32(datTogRef.Position[1])
			enc.Float32(datTogRef.Position[2])
			enc.Float32(datTogRef.Rotation[0])
			enc.Float32(datTogRef.Rotation[1])
			enc.Float32(datTogRef.Rotation[2])
			enc.Float32(datTogRef.Scale[0])
			enc.Float32(datTogRef.Scale[1])
			enc.Float32(datTogRef.Scale[2])
			enc.Float32(datTogRef.Adjust)
		}
	}

	return nil
}
