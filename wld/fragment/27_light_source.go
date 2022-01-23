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
	hashIndex    uint32
}

func LoadLightSource(r io.ReadSeeker) (common.WldFragmenter, error) {
	e := &LightSource{}
	err := parseLightSource(r, e)
	if err != nil {
		return nil, fmt.Errorf("parse LightSource: %w", err)
	}
	return e, nil
}

func parseLightSource(r io.ReadSeeker, e *LightSource) error {
	if e == nil {
		return fmt.Errorf("lightsource is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &e.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	if value&1 == 1 {
		e.IsPlacedLightSource = true
	}
	if value&4 == 4 {
		e.IsColoredLight = true
	}

	if !e.IsPlacedLightSource {
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
	if e.IsColoredLight {
		err = binary.Read(r, binary.LittleEndian, &e.Attentuation)
		if err != nil {
			return fmt.Errorf("read attentuation: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &e.Color.A)
		if err != nil {
			return fmt.Errorf("read color alpha: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &e.Color.R)
		if err != nil {
			return fmt.Errorf("read color red: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &e.Color.G)
		if err != nil {
			return fmt.Errorf("read color green: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &e.Color.B)
		if err != nil {
			return fmt.Errorf("read color blue: %w", err)
		}
		return nil
	}
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown noncolored: %w", err)
	}
	e.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	return nil
}

func (e *LightSource) FragmentType() string {
	return "Light Source"
}

func (e *LightSource) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
