package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/xackery/quail/common"
)

// LightSource information
type LightSource struct {
	// IsPlacedLightSource if used in the light.wld and is not if used in the main zone file
	IsPlacedLightSource bool
	// IsColoredLight returns true and impacts fragment size
	IsColoredLight bool
	// Color of the light, if applicable
	Color color.RGBA
	//Attenuation (?) - guess from Windcatcher. Not sure what it is.
	Attentuation uint32
	name         string
}

func LoadLightSource(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &LightSource{}
	err := parseLightSource(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse LightSource: %w", err)
	}
	return e, nil
}

func parseLightSource(r io.ReadSeeker, v *LightSource) error {
	if v == nil {
		return fmt.Errorf("lightsource is nil")
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

	if value&1 == 1 {
		v.IsPlacedLightSource = true
	}
	if value&4 == 4 {
		v.IsColoredLight = true
	}

	if !v.IsPlacedLightSource {
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read unknown: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("read unknown6: %w", err)
		}
		return nil
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown1: %w", err)
	}
	if v.IsColoredLight {
		err = binary.Read(r, binary.LittleEndian, &v.Attentuation)
		if err != nil {
			return fmt.Errorf("read attentuation: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &v.Color.A)
		if err != nil {
			return fmt.Errorf("read color alpha: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &v.Color.R)
		if err != nil {
			return fmt.Errorf("read color red: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &v.Color.G)
		if err != nil {
			return fmt.Errorf("read color green: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &v.Color.B)
		if err != nil {
			return fmt.Errorf("read color blue: %w", err)
		}
		return nil
	}
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown noncolored: %w", err)
	}
	v.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	return nil
}

func (e *LightSource) FragmentType() string {
	return "Light Source"
}

func (e *LightSource) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
