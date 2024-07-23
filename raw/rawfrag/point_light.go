package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragPointLight is PointLight in libeq, Light Info in openzone, POINTLIGHT in wld, LightInstance in lantern
type WldFragPointLight struct {
	NameRef  int32      `yaml:"name_ref"`
	LightRef int32      `yaml:"light_ref"`
	Flags    uint32     `yaml:"flags"`
	Location [3]float32 `yaml:"location"`
	Radius   float32    `yaml:"radius"`
}

func (e *WldFragPointLight) FragCode() int {
	return FragCodePointLight
}

func (e *WldFragPointLight) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.LightRef)
	enc.Uint32(e.Flags)

	enc.Float32(e.Location[0])
	enc.Float32(e.Location[1])
	enc.Float32(e.Location[2])
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
	e.Location[0] = dec.Float32()
	e.Location[1] = dec.Float32()
	e.Location[2] = dec.Float32()
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
