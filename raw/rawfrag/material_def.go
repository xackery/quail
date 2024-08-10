package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragMaterialDef is MaterialDef in libeq, Texture in openzone, MATERIALDEFINITION in wld, Material in lantern
type WldFragMaterialDef struct {
	NameRef         int32    `yaml:"name_ref"`
	Flags           uint32   `yaml:"flags"`
	RenderMethod    uint32   `yaml:"render_method"`
	RGBPen          [4]uint8 `yaml:"rgb_pen"`
	Brightness      float32  `yaml:"brightness"`
	ScaledAmbient   float32  `yaml:"scaled_ambient"`
	SimpleSpriteRef uint32   `yaml:"sprite_instance_ref"`
	Pair1           uint32
	Pair2           float32
}

func (e *WldFragMaterialDef) FragCode() int {
	return FragCodeMaterialDef
}

func (e *WldFragMaterialDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.RenderMethod)
	enc.Uint8(e.RGBPen[0])
	enc.Uint8(e.RGBPen[1])
	enc.Uint8(e.RGBPen[2])
	enc.Uint8(e.RGBPen[3])
	enc.Float32(e.Brightness)
	enc.Float32(e.ScaledAmbient)
	enc.Uint32(e.SimpleSpriteRef)
	if e.Flags&0x2 != 0 {
		enc.Uint32(e.Pair1)
		enc.Float32(e.Pair2)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragMaterialDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.RenderMethod = dec.Uint32()
	e.RGBPen = [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()}
	e.Brightness = dec.Float32()
	e.ScaledAmbient = dec.Float32()
	e.SimpleSpriteRef = dec.Uint32()
	if e.Flags&0x2 != 0 {
		e.Pair1 = dec.Uint32()
		e.Pair2 = dec.Float32()
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
