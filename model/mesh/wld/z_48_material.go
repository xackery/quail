package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ghostiam/binstruct"
	"github.com/xackery/quail/log"
)

// material 0x30 48
type material struct {
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
	pairs         [2]uint32 `bin:"-"` //This only exists if bit 1 of flags is set. Both fields usually contain 0.
}

func (e *WLD) materialRead(r io.ReadSeeker, fragmentOffset int) error {
	def := &material{}

	dec := binstruct.NewDecoder(r, binary.LittleEndian)
	err := dec.Decode(def)
	if err != nil {
		return fmt.Errorf("decode material: %w", err)
	}

	if def.Flags&0x1 != 0 {
		err = binary.Read(r, binary.LittleEndian, &def.pairs)
		if err != nil {
			return fmt.Errorf("read pairs: %w", err)
		}
	}

	log.Debugf("material: %+v\n", def)
	log.Debugf("%+v", def)
	e.fragments[fragmentOffset] = def
	return nil
}

func (v *material) build(e *WLD) error {
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
