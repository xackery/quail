package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/xackery/quail/pfs/archive"
)

// MaterialDef information
type MaterialDef struct {
	name string
	//BitmapInfoReference
	// shaderType is the way to render the material
	shaderType int
	// materialType is also part of rendering material
	materialType int
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

func LoadMaterialDef(r io.ReadSeeker) (archive.WldFragmenter, error) {
	m := &MaterialDef{}
	err := parseMaterialDef(r, m)
	if err != nil {
		return nil, fmt.Errorf("parse MaterialDef: %w", err)
	}
	return m, nil
}

func parseMaterialDef(r io.ReadSeeker, v *MaterialDef) error {
	if v == nil {
		return fmt.Errorf("MaterialDef is nil")
	}
	var value uint32

	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHashIndex: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	var params int32
	err = binary.Read(r, binary.LittleEndian, &params)
	if err != nil {
		return fmt.Errorf("read params: %w", err)
	}

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

	//if value != 0 {
	//BitmapInfoReference = fragments[fragmentReference - 1] as BitmapInfoReference;
	//}

	v.materialType = int(int64(params) & ^0x80000000)
	switch v.materialType {
	case MaterialTypeBoundary:
		v.shaderType = ShaderTypeBoundary
	case MaterialTypeInvisibleUnknown, MaterialTypeInvisibleUnknown2, MaterialTypeInvisibleUnknown3:
		v.shaderType = ShaderTypeInvisible
	case MaterialTypeDiffuse, MaterialTypeDiffuse2, MaterialTypeDiffuse3, MaterialTypeDiffuse4, MaterialTypeDiffuse6, MaterialTypeDiffuse7, MaterialTypeDiffuse8, MaterialTypeCompleteUnknown, MaterialTypeTransparentMaskedPassable:
		v.shaderType = ShaderTypeDiffuse
	case MaterialTypeTransparent25:
		v.shaderType = ShaderTypeTransparent25
	case MaterialTypeTransparent50:
		v.shaderType = ShaderTypeTransparent50
	case MaterialTypeTransparent75:
		v.shaderType = ShaderTypeTransparent75
	case MaterialTypeTransparentAdditive:
		v.shaderType = ShaderTypeTransparentAdditive
	case MaterialTypeTransparentAdditiveUnlit:
		v.shaderType = ShaderTypeTransparentAdditiveUnlit
	case MaterialTypeTransparentMasked, MaterialTypeDiffuse5:
		v.shaderType = ShaderTypeTransparentMasked
	case MaterialTypeDiffuseSkydome:
		v.shaderType = ShaderTypeDiffuseSkydome
	case MaterialTypeTransparentSkydome:
		v.shaderType = ShaderTypeTransparentSkydome
	case MaterialTypeTransparentAdditiveUnlitSkydome:
		v.shaderType = ShaderTypeTransparentAdditiveUnlitSkydome
	default:
		//m.ShaderType = BitmapInfoReference == null ? ShaderTypeInvisible : ShaderTypeDiffuse;
	}
	return nil
}

func (m *MaterialDef) FragmentType() string {
	return "MaterialDef"
}

func (e *MaterialDef) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}

func (e *MaterialDef) Name() string {
	return e.name
}

func (e *MaterialDef) ShaderType() int {
	return e.shaderType
}

func (e *MaterialDef) MaterialType() int {
	return e.materialType
}
