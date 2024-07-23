package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragLightDef is LightDef in libeq, WldFragLightDef Source in openzone, LIGHT (ref) in wld, LightSource in lantern
type WldFragLightDef struct {
	NameRef         int32        `yaml:"name_ref"`
	Flags           uint32       `yaml:"flags"`
	FrameCurrentRef uint32       `yaml:"frame_current_ref"`
	Sleep           uint32       `yaml:"sleep"`
	LightLevels     []float32    `yaml:"light_levels"`
	Colors          [][3]float32 `yaml:"colors"`
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
			enc.Float32(color[0])
			enc.Float32(color[1])
			enc.Float32(color[2])
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
			e.Colors = append(e.Colors, [3]float32{dec.Float32(), dec.Float32(), dec.Float32()})
		}
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
