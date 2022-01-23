package fragment

import (
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
	hashIndex    uint32
}

func LoadLightSource(r io.ReadSeeker) (common.WldFragmenter, error) {
	l := &LightSource{}
	err := parseLightSource(r, l)
	if err != nil {
		return nil, fmt.Errorf("parse LightSource: %w", err)
	}
	return l, nil
}

func parseLightSource(r io.ReadSeeker, l *LightSource) error {
	if l == nil {
		return fmt.Errorf("lightsource is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &l.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	if value&1 == 1 {
		l.IsPlacedLightSource = true
	}
	if value&4 == 4 {
		l.IsColoredLight = true
	}

	if !l.IsPlacedLightSource {
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
	if l.IsColoredLight {
		err = binary.Read(r, binary.LittleEndian, &l.Attentuation)
		if err != nil {
			return fmt.Errorf("read attentuation: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &l.Color.A)
		if err != nil {
			return fmt.Errorf("read color alpha: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &l.Color.R)
		if err != nil {
			return fmt.Errorf("read color red: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &l.Color.G)
		if err != nil {
			return fmt.Errorf("read color green: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &l.Color.B)
		if err != nil {
			return fmt.Errorf("read color blue: %w", err)
		}
		return nil
	}
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown noncolored: %w", err)
	}
	l.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	return nil
}

func (l *LightSource) FragmentType() string {
	return "Light Source"
}
