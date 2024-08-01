package rawfrag

import "github.com/xackery/encdec"

type RenderInfo struct {
	Method                uint32
	Flags                 uint8
	Pen                   uint32
	Brightness            float32
	ScaledAmbient         float32
	SimpleSpriteReference uint32
	UVInfoOrigin          [3]float32
	UVInfoUAxis           [3]float32
	UVInfoVAxis           [3]float32
	Uvs                   [][2]float32
}

func (e *RenderInfo) Write(enc *encdec.Encoder) error {
	enc.Uint32(e.Method)
	enc.Uint8(e.Flags)
	enc.Uint32(e.Pen)
	enc.Float32(e.Brightness)
	enc.Float32(e.ScaledAmbient)
	enc.Uint32(e.SimpleSpriteReference)
	enc.Float32(e.UVInfoOrigin[0])
	enc.Float32(e.UVInfoOrigin[1])
	enc.Float32(e.UVInfoOrigin[2])
	enc.Float32(e.UVInfoUAxis[0])
	enc.Float32(e.UVInfoUAxis[1])
	enc.Float32(e.UVInfoUAxis[2])
	enc.Float32(e.UVInfoVAxis[0])
	enc.Float32(e.UVInfoVAxis[1])
	enc.Float32(e.UVInfoVAxis[2])
	enc.Uint32(uint32(len(e.Uvs)))
	for _, uv := range e.Uvs {
		enc.Float32(uv[0])
		enc.Float32(uv[1])
	}
	return nil

}
