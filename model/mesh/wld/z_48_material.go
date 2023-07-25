package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
)

// Material 0x30 48
type Material struct {
	NameRef int32
	Flags   uint32
	/* // Bit 0 ........ Apparently must be 1 if the texture isn’t transparent.
	   // Bit 1 ........ Set to 1 if the texture is masked (e.g. tree leaves).
	   // Bit 2 ........ Set to 1 if the texture is semi-transparent but not masked.
	   // Bit 3 ........ Set to 1 if the texture is masked and semi-transparent.
	   // Bit 4 ........ Set to 1 if the texture is masked but not semi-transparent.
	   // Bit 31 ...... Apparently must be 1 if the texture isn’t transparent.
	*/
	RenderMethod  uint32
	RGBPen        uint32 // This typically contains 0x004E4E4E but has also been known to contain 0xB2B2B2.
	Brightness    float32
	ScaledAmbient float32
	TextureRef    uint32
	Pairs         [2]uint32 //This only exists if bit 1 of flags is set. Both fields usually contain 0.
}

func (e *WLD) materialRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &Material{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	def.NameRef = dec.Int32()
	def.Flags = dec.Uint32()
	def.RenderMethod = dec.Uint32()
	def.RGBPen = dec.Uint32()
	def.Brightness = dec.Float32()
	def.ScaledAmbient = dec.Float32()
	def.TextureRef = dec.Uint32()
	if def.Flags&0x1 != 0 {
		def.Pairs[0] = dec.Uint32()
		def.Pairs[1] = dec.Uint32()
	}

	if dec.Error() != nil {
		return fmt.Errorf("materialRead: %v", dec.Error())
	}

	log.Debugf("%+v", def)
	e.Fragments[fragmentOffset] = def
	return nil
}

func (v *Material) build(e *WLD) error {
	return nil
}

const (
	// Used for boundaries that are not rendered. TextInfoReference can be null or have reference.
	MaterialTypeBoundary = 0x0
	// Standard diffuse shader
	MaterialTypeDiffuse = 0x01
	// Diffuse variant
	MaterialTypeDiffuse2 = 0x02
	// Transparent with 0.5 blend strength
	MaterialTypeTransparent50 = 0x05
	// Transparent with 0.25 blend strength
	MaterialTypeTransparent25 = 0x09
	// Transparent with 0.75 blend strength
	MaterialTypeTransparent75 = 0x0A
	// Non solid surfaces that shouldn't really be masked
	MaterialTypeTransparentMaskedPassable       = 0x07
	MaterialTypeTransparentAdditiveUnlit        = 0x0B
	MaterialTypeTransparentMasked               = 0x13
	MaterialTypeDiffuse3                        = 0x14
	MaterialTypeDiffuse4                        = 0x15
	MaterialTypeTransparentAdditive             = 0x17
	MaterialTypeDiffuse5                        = 0x19
	MaterialTypeInvisibleUnknown                = 0x53
	MaterialTypeDiffuse6                        = 0x553
	MaterialTypeCompleteUnknown                 = 0x1A
	MaterialTypeDiffuse7                        = 0x12
	MaterialTypeDiffuse8                        = 0x31
	MaterialTypeInvisibleUnknown2               = 0x4B
	MaterialTypeDiffuseSkydome                  = 0x0D // Need to confirm
	MaterialTypeTransparentSkydome              = 0x0F // Need to confirm
	MaterialTypeTransparentAdditiveUnlitSkydome = 0x10
	MaterialTypeInvisibleUnknown3               = 0x03
)

func (e *WLD) materialWrite(w io.Writer, fragmentOffset int) error {
	return fmt.Errorf("not implemented")
}
