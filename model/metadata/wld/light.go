package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// Light is LightDef in libeq, Light Source in openzone, LIGHT (ref) in wld, LightSource in lantern
type Light struct {
	FragName        string           `yaml:"frag_name"`
	NameRef         int32            `yaml:"name_ref"`
	Flags           uint32           `yaml:"flags"`
	FrameCurrentRef uint32           `yaml:"frame_current_ref"`
	Sleep           uint32           `yaml:"sleep"`
	LightLevels     []float32        `yaml:"light_levels"`
	Colors          []common.Vector3 `yaml:"colors"`
}

func (e *Light) FragCode() int {
	return 0x1B
}

func (e *Light) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	frameCount := uint32(len(e.LightLevels))
	enc.Uint32(frameCount)
	if e.Flags&0x1 != 0 {
		enc.Uint32(e.FrameCurrentRef)
	}
	if e.Flags&0x2 != 0 {
		enc.Uint32(e.Sleep)
	}
	if e.Flags&0x4 != 0 {
		for _, lightLevel := range e.LightLevels {
			enc.Float32(lightLevel)
		}
	}
	if e.Flags&0x10 != 0 {
		for _, color := range e.Colors {
			enc.Float32(color.X)
			enc.Float32(color.Y)
			enc.Float32(color.Z)
		}
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeLight(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Light{}
	d.FragName = common.FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	frameCount := dec.Uint32()
	if d.Flags&0x1 != 0 {
		d.FrameCurrentRef = dec.Uint32()
	}
	if d.Flags&0x2 != 0 {
		d.Sleep = dec.Uint32()
	}
	if d.Flags&0x4 != 0 {
		for i := uint32(0); i < frameCount; i++ {
			d.LightLevels = append(d.LightLevels, dec.Float32())
		}
	}
	if d.Flags&0x10 != 0 {
		for i := uint32(0); i < frameCount; i++ {
			var color common.Vector3
			color.X = dec.Float32()
			color.Y = dec.Float32()
			color.Z = dec.Float32()
			d.Colors = append(d.Colors, color)
		}
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// LightRef is Light in libeq, Light Source Reference in openzone, POINTLIGHTT ?? in wld, LightSourceReference in lantern
type LightRef struct {
	FragName    string `yaml:"frag_name"`
	NameRef     int32  `yaml:"name_ref"`
	LightDefRef int32  `yaml:"light_def_ref"`
	Flags       uint32 `yaml:"flags"`
}

func (e *LightRef) FragCode() int {
	return 0x1C
}

func (e *LightRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightDefRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeLightRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &LightRef{}
	d.FragName = common.FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.LightDefRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// PointLightOld is empty in libeq, empty in openzone, POINTLIGHT?? in wld
type PointLightOld struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	Flags    uint32 `yaml:"flags"`
}

func (e *PointLightOld) FragCode() int {
	return 0x1D
}

func (e *PointLightOld) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodePointLightOld(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &PointLightOld{}
	d.FragName = common.FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// PointLightOldRef is empty in libeq, empty in openzone, empty in wld
type PointLightOldRef struct {
	FragName      string `yaml:"frag_name"`
	NameRef       int32  `yaml:"name_ref"`
	PointLightRef int32  `yaml:"point_light_ref"`
}

func (e *PointLightOldRef) FragCode() int {
	return 0x1E
}

func (e *PointLightOldRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.PointLightRef)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodePointLightOldRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &PointLightOldRef{}
	d.FragName = common.FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.PointLightRef = dec.Int32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// DirectionalLigtOld is empty in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type DirectionalLightOld struct {
	FragName string `yaml:"frag_name"`
}

func (e *DirectionalLightOld) FragCode() int {
	return 0x25
}

func (e *DirectionalLightOld) Encode(w io.Writer) error {
	return nil
}

func decodeDirectionalLightOld(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &DirectionalLightOld{}
	d.FragName = common.FragName(d.FragCode())
	return d, nil
}

// PointLight is PointLight in libeq, Light Info in openzone, POINTLIGHT in wld, LightInstance in lantern
type PointLight struct {
	FragName string `yaml:"frag_name"`
}

func (e *PointLight) FragCode() int {
	return 0x28
}

func (e *PointLight) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func decodePointLight(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &PointLight{}
	d.FragName = common.FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// AmbientLight is AmbientLight in libeq, Ambient Light in openzone, AMBIENTLIGHT in wld, AmbientLight in lantern
type AmbientLight struct {
	FragName string `yaml:"frag_name"`
}

func (e *AmbientLight) FragCode() int {
	return 0x2A
}

func (e *AmbientLight) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func decodeAmbientLight(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &AmbientLight{}
	d.FragName = common.FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// DirectionalLight is DirectionalLight in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type DirectionalLight struct {
	FragName string `yaml:"frag_name"`
}

func (e *DirectionalLight) FragCode() int {
	return 0x2B
}

func (e *DirectionalLight) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func decodeDirectionalLight(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &DirectionalLight{}
	d.FragName = common.FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}
