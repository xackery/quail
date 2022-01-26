package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/xackery/quail/common"
)

// Material information
type Material struct {
	name string
	//BitmapInfoReference
	// ShaderType is the way to render the material
	ShaderType int
	// MaterialType is also part of rendering material
	MaterialType int
	// IsHandled is used when an alternative character skin is needed
	IsHandled bool
}

const (
	ShaderTypeDiffuse                         = 0
	ShaderTypeTransparent25                   = 1
	ShaderTypeTransparent50                   = 2
	ShaderTypeTransparent75                   = 3
	ShaderTypeTransparentAdditive             = 4
	ShaderTypeTransparentAdditiveUnlit        = 5
	ShaderTypeTransparentMasked               = 6
	ShaderTypeDiffuseSkydome                  = 7
	ShaderTypeTransparentSkydome              = 8
	ShaderTypeTransparentAdditiveUnlitSkydome = 9
	ShaderTypeInvisible                       = 10
	ShaderTypeBoundary                        = 11
)

func LoadMaterial(r io.ReadSeeker) (common.WldFragmenter, error) {
	m := &Material{}
	err := parseMaterial(r, m)
	if err != nil {
		return nil, fmt.Errorf("parse Material: %w", err)
	}
	return m, nil
}

func parseMaterial(r io.ReadSeeker, v *Material) error {
	if v == nil {
		return fmt.Errorf("Material is nil")
	}
	var value uint32

	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}
	//TODO: flags support
	var params int32
	err = binary.Read(r, binary.LittleEndian, &params)
	if err != nil {
		return fmt.Errorf("read params: %w", err)
	}

	//TODO: figure out color
	rgba := color.RGBA{}
	err = binary.Read(r, binary.LittleEndian, &rgba.R)
	if err != nil {
		return fmt.Errorf("read color red: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &rgba.G)
	if err != nil {
		return fmt.Errorf("read color green: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &rgba.B)
	if err != nil {
		return fmt.Errorf("read color blue: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &rgba.A)
	if err != nil {
		return fmt.Errorf("read color alpha: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &rgba.A)
	if err != nil {
		return fmt.Errorf("read unknown float 1: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &rgba.A)
	if err != nil {
		return fmt.Errorf("read unknown float 2: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read fragment reference: %w", err)
	}
	// TODO: add bitmapinforeference support
	//if value != 0 {
	//BitmapInfoReference = fragments[fragmentReference - 1] as BitmapInfoReference;
	//}

	v.MaterialType = int(int64(params) & ^0x80000000)
	switch v.MaterialType {
	case MaterialTypeBoundary:
		v.ShaderType = ShaderTypeBoundary
	case MaterialTypeInvisibleUnknown, MaterialTypeInvisibleUnknown2, MaterialTypeInvisibleUnknown3:
		v.ShaderType = ShaderTypeInvisible
	case MaterialTypeDiffuse, MaterialTypeDiffuse2, MaterialTypeDiffuse3, MaterialTypeDiffuse4, MaterialTypeDiffuse6, MaterialTypeDiffuse7, MaterialTypeDiffuse8, MaterialTypeCompleteUnknown, MaterialTypeTransparentMaskedPassable:
		v.ShaderType = ShaderTypeDiffuse
	case MaterialTypeTransparent25:
		v.ShaderType = ShaderTypeTransparent25
	case MaterialTypeTransparent50:
		v.ShaderType = ShaderTypeTransparent50
	case MaterialTypeTransparent75:
		v.ShaderType = ShaderTypeTransparent75
	case MaterialTypeTransparentAdditive:
		v.ShaderType = ShaderTypeTransparentAdditive
	case MaterialTypeTransparentAdditiveUnlit:
		v.ShaderType = ShaderTypeTransparentAdditiveUnlit
	case MaterialTypeTransparentMasked, MaterialTypeDiffuse5:
		v.ShaderType = ShaderTypeTransparentMasked
	case MaterialTypeDiffuseSkydome:
		v.ShaderType = ShaderTypeDiffuseSkydome
	case MaterialTypeTransparentSkydome:
		v.ShaderType = ShaderTypeTransparentSkydome
	case MaterialTypeTransparentAdditiveUnlitSkydome:
		v.ShaderType = ShaderTypeTransparentAdditiveUnlitSkydome
	default:
		//m.ShaderType = BitmapInfoReference == null ? ShaderTypeInvisible : ShaderTypeDiffuse;
	}
	return nil
}

func (m *Material) FragmentType() string {
	return "Material"
}

func (e *Material) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
