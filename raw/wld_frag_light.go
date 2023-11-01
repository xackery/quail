package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragLight is LightDef in libeq, WldFragLight Source in openzone, LIGHT (ref) in wld, LightSource in lantern
type WldFragLight struct {
	FragName        string    `yaml:"frag_name"`
	NameRef         int32     `yaml:"name_ref"`
	Flags           uint32    `yaml:"flags"`
	FrameCurrentRef uint32    `yaml:"frame_current_ref"`
	Sleep           uint32    `yaml:"sleep"`
	LightLevels     []float32 `yaml:"light_levels"`
	Colors          []Vector3 `yaml:"colors"`
}

func (e *WldFragLight) FragCode() int {
	return 0x1B
}

func (e *WldFragLight) Encode(w io.Writer) error {
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

func decodeLight(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragLight{}
	d.FragName = FragName(d.FragCode())
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
			var color Vector3
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

// WldFragLightRef is Light in libeq, Light Source Reference in openzone, POINTLIGHTT ?? in wld, LightSourceReference in lantern
type WldFragLightRef struct {
	FragName    string `yaml:"frag_name"`
	NameRef     int32  `yaml:"name_ref"`
	LightDefRef int32  `yaml:"light_def_ref"`
	Flags       uint32 `yaml:"flags"`
}

func (e *WldFragLightRef) FragCode() int {
	return 0x1C
}

func (e *WldFragLightRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightDefRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeLightRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragLightRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.LightDefRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragPointLightOld is empty in libeq, empty in openzone, POINTLIGHT?? in wld
type WldFragPointLightOld struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	Flags    uint32 `yaml:"flags"`
}

func (e *WldFragPointLightOld) FragCode() int {
	return 0x1D
}

func (e *WldFragPointLightOld) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodePointLightOld(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragPointLightOld{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragPointLightOldRef is empty in libeq, empty in openzone, empty in wld
type WldFragPointLightOldRef struct {
	FragName      string `yaml:"frag_name"`
	NameRef       int32  `yaml:"name_ref"`
	PointLightRef int32  `yaml:"point_light_ref"`
}

func (e *WldFragPointLightOldRef) FragCode() int {
	return 0x1E
}

func (e *WldFragPointLightOldRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.PointLightRef)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodePointLightOldRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragPointLightOldRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.PointLightRef = dec.Int32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// DirectionalLigtOld is empty in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLightOld struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragDirectionalLightOld) FragCode() int {
	return 0x25
}

func (e *WldFragDirectionalLightOld) Encode(w io.Writer) error {
	return nil
}

func decodeDirectionalLightOld(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragDirectionalLightOld{}
	d.FragName = FragName(d.FragCode())
	return d, nil
}

// WldFragPointLight is WldFragPointLight in libeq, Light Info in openzone, POINTLIGHT in wld, LightInstance in lantern
type WldFragPointLight struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragPointLight) FragCode() int {
	return 0x28
}

func (e *WldFragPointLight) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func decodePointLight(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragPointLight{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragAmbientLight is WldFragAmbientLight in libeq, Ambient Light in openzone, AMBIENTLIGHT in wld, WldFragAmbientLight in lantern
type WldFragAmbientLight struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragAmbientLight) FragCode() int {
	return 0x2A
}

func (e *WldFragAmbientLight) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func decodeAmbientLight(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragAmbientLight{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragDirectionalLight is WldFragDirectionalLight in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLight struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragDirectionalLight) FragCode() int {
	return 0x2B
}

func (e *WldFragDirectionalLight) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func decodeDirectionalLight(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragDirectionalLight{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}
