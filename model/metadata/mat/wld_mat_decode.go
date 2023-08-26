// Package mat repesents metadata related to material decoding
package mat

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// DecodeTextureImages decodes texture images
func DecodeTextureImages(material *common.Material, nameRef *int32, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	*nameRef = dec.Int32()
	textureCount := dec.Int32()

	for i := 0; i < int(textureCount+1); i++ {
		nameLength := dec.Uint16()
		name := helper.ReadStringHash(dec.Bytes(int(nameLength)))
		if i != 0 {
			material.Animation.Textures = append(material.Animation.Textures, name)
			continue
		}

		prop := &common.MaterialProperty{
			Name:     "e_TextureDiffuse0",
			Category: 2,
			Value:    name,
		}
		material.Properties = append(material.Properties, prop)
	}
	if dec.Error() != nil {
		return fmt.Errorf("decodeTextureImages: %s", dec.Error())
	}
	return nil
}

// DecodeTexture decodes texture
func DecodeTexture(material *common.Material, nameRef *int32, textureRefs []*int32, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	*nameRef = dec.Int32()
	flags := dec.Uint32()
	textureCount := dec.Int32()
	if flags&0x20 != 0 {
		dec.Uint32() // textureCurrent
	}
	if flags&0x08 != 0 && flags&0x10 != 0 {
		material.Animation.Sleep = dec.Uint32()
	}
	for i := 0; i < int(textureCount); i++ {
		ref := dec.Int32()
		textureRefs = append(textureRefs, &ref)
	}

	if dec.Error() != nil {
		return fmt.Errorf("decodeTexture: %s", dec.Error())
	}
	return nil
}

func DecodeMaterialDef(material *common.Material, nameRef *int32, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	*nameRef = dec.Int32()
	flags := dec.Uint32()
	renderMethod := dec.Uint32() &^ 0x80000000

	/*
	   // Most flags are _unknown_, however:
	   // Bit 0 ........ Apparently must be 1 if the texture isn’t transparent.
	   // Bit 1 ........ Set to 1 if the texture is masked (e.g. tree leaves).
	   // Bit 2 ........ Set to 1 if the texture is semi-transparent but not masked.
	   // Bit 3 ........ Set to 1 if the texture is masked and semi-transparent.
	   // Bit 4 ........ Set to 1 if the texture is masked but not semi-transparent.
	   // Bit 31 ...... Apparently must be 1 if the texture isn’t transparent.
	*/

	switch renderMethod {
	// Standard diffuse shader
	case 0x00:
		material.ShaderName = "Boundary"
	case 0x01: // Diffuse
	case 0x12:
	case 0x31:
	case 0x14:
	case 0x15:
	case 0x02: // Diffuse Variant
	case 0x19:
	case 0x553:
	case 0x07: // Non solid surfaces that shouldn't really be masked
		//material.ShaderName = "Diffuse2"
	case 0x05: // Transparent with 0.5 blend strength
		material.ShaderName = "Transparent50"
	case 0x09: // Transparent with 0.25 blend strength
		material.ShaderName = "Transparent25"
	case 0x0A: // Transparent with 0.75 blend strength
		material.ShaderName = "Transparent75"
	case 0x0B:
		material.ShaderName = "TransparentAdditiveUnlit"
	case 0x13:
		material.ShaderName = "TransparentMasked"
	case 0x17:
		material.ShaderName = "TransparentAdditive"
	case 0x53:
	case 0x4B:
	case 0x03:
		material.ShaderName = "invis"
	case 0x1A: // TODO: Analyze thi:
		material.ShaderName = "CompleteUnknown"
	case 0x0D:
		material.ShaderName = "DiffuseSkydome"
	case 0x0F:
		material.ShaderName = "TransparentSkydome"
	case 0x10:
		material.ShaderName = "TransparentAdditiveUnlitSkydome"
	}
	// This typically contains 0x004E4E4E but has also bee known to contain 0xB2B2B2.
	// Could this be an RGB reflectivity value?
	dec.Uint32() // RGBPEN %d, %d, %d

	dec.Float32() // BRIGHTNESS %f

	dec.Uint32() // textureRef, ignored for now

	if flags&0x2 == 0x2 {
		dec.Uint32()
		dec.Uint32()
	}

	if dec.Error() != nil {
		return fmt.Errorf("decodeMaterial: %s", dec.Error())
	}
	return nil
}
