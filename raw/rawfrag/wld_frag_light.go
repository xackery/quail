package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragLightDef is LightDef in libeq, WldFragLightDef Source in openzone, LIGHT (ref) in wld, LightSource in lantern
type WldFragLightDef struct {
	NameRef         int32           `yaml:"name_ref"`
	Flags           uint32          `yaml:"flags"`
	FrameCurrentRef uint32          `yaml:"frame_current_ref"`
	Sleep           uint32          `yaml:"sleep"`
	LightLevels     []float32       `yaml:"light_levels"`
	Colors          []model.Vector3 `yaml:"colors"`
}

func (e *WldFragLightDef) FragCode() int {
	return FragCodeLightDef
}

func (e *WldFragLightDef) Write(w io.Writer) error {
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
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragLightDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	frameCount := dec.Uint32()
	if e.Flags&0x1 != 0 {
		e.FrameCurrentRef = dec.Uint32()
	}
	if e.Flags&0x2 != 0 {
		e.Sleep = dec.Uint32()
	}
	if e.Flags&0x4 != 0 {
		for i := uint32(0); i < frameCount; i++ {
			e.LightLevels = append(e.LightLevels, dec.Float32())
		}
	}
	if e.Flags&0x10 != 0 {
		for i := uint32(0); i < frameCount; i++ {
			var color model.Vector3
			color.X = dec.Float32()
			color.Y = dec.Float32()
			color.Z = dec.Float32()
			e.Colors = append(e.Colors, color)
		}
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragLight is Light in libeq, Light Source Reference in openzone, POINTLIGHTT ?? in wld, LightSourceReference in lantern
type WldFragLight struct {
	NameRef     int32  `yaml:"name_ref"`
	LightDefRef int32  `yaml:"light_def_ref"`
	Flags       uint32 `yaml:"flags"`
}

func (e *WldFragLight) FragCode() int {
	return FragCodeLight
}

func (e *WldFragLight) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightDefRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragLight) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.LightDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragPointLightOld is empty in libeq, empty in openzone, POINTLIGHT?? in wld
type WldFragPointLightOld struct {
	NameRef int32  `yaml:"name_ref"`
	Flags   uint32 `yaml:"flags"`
}

func (e *WldFragPointLightOld) FragCode() int {
	return FragCodePointLightOld
}

func (e *WldFragPointLightOld) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPointLightOld) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragPointLightOldDef is empty in libeq, empty in openzone, empty in wld
type WldFragPointLightOldDef struct {
	NameRef       int32 `yaml:"name_ref"`
	PointLightRef int32 `yaml:"point_light_ref"`
}

func (e *WldFragPointLightOldDef) FragCode() int {
	return FragCodePointLightOldDef
}

func (e *WldFragPointLightOldDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.PointLightRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPointLightOldDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.PointLightRef = dec.Int32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// DirectionalLigtOld is empty in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLightOld struct {
}

func (e *WldFragDirectionalLightOld) FragCode() int {
	return FragCodeDirectionalLightOld
}

func (e *WldFragDirectionalLightOld) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDirectionalLightOld) Read(r io.ReadSeeker) error {
	return nil
}

// WldFragPointLight is PointLight in libeq, Light Info in openzone, POINTLIGHT in wld, LightInstance in lantern
type WldFragPointLight struct {
	NameRef  int32   `yaml:"name_ref"`
	LightRef int32   `yaml:"light_ref"`
	Flags    uint32  `yaml:"flags"`
	X        float32 `yaml:"x"`
	Y        float32 `yaml:"y"`
	Z        float32 `yaml:"z"`
	Radius   float32 `yaml:"radius"`
}

func (e *WldFragPointLight) FragCode() int {
	return FragCodePointLight
}

func (e *WldFragPointLight) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightRef)
	enc.Uint32(e.Flags)
	enc.Float32(e.X)
	enc.Float32(e.Y)
	enc.Float32(e.Z)
	enc.Float32(e.Radius)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragPointLight) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.LightRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.X = dec.Float32()
	e.Y = dec.Float32()
	e.Z = dec.Float32()
	e.Radius = dec.Float32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragAmbientLight is AmbientLight in libeq, Ambient Light in openzone, AMBIENTLIGHT in wld, WldFragAmbientLight in lantern
type WldFragAmbientLight struct {
	NameRef  int32    `yaml:"name_ref"`
	LightRef int32    `yaml:"light_ref"`
	Flags    uint32   `yaml:"flags"`
	Regions  []uint32 `yaml:"regions"`
}

func (e *WldFragAmbientLight) FragCode() int {
	return FragCodeAmbientLight
}

func (e *WldFragAmbientLight) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Regions)))
	for _, region := range e.Regions {
		enc.Uint32(region)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragAmbientLight) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.LightRef = dec.Int32()
	e.Flags = dec.Uint32()
	regionCount := dec.Uint32()
	for i := uint32(0); i < regionCount; i++ {
		e.Regions = append(e.Regions, dec.Uint32())
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragDirectionalLight is DirectionalLight in libeq, empty in openzone, DIRECTIONALLIGHT in wld
type WldFragDirectionalLight struct {
}

func (e *WldFragDirectionalLight) FragCode() int {
	return FragCodeDirectionalLight
}

func (e *WldFragDirectionalLight) Write(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func (e *WldFragDirectionalLight) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
