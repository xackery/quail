package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragSprite4DDef is Sprite4DDef in libeq, empty in openzone, 4DSPRITEDEF in wld
type WldFragSprite4DDef struct {
	NameRef         int32         `yaml:"name_ref"`
	Flags           uint32        `yaml:"flags"`
	PolyRef         int32         `yaml:"poly_ref"`
	CenterOffset    model.Vector3 `yaml:"center_offset"`
	Radius          float32       `yaml:"radius"`
	CurrentFrame    uint32        `yaml:"current_frame"`
	Sleep           uint32        `yaml:"sleep"`
	SpriteFragments []uint32      `yaml:"sprite_fragments"`
}

func (e *WldFragSprite4DDef) FragCode() int {
	return FragCodeSprite4DDef
}

func (e *WldFragSprite4DDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.SpriteFragments)))
	enc.Int32(e.PolyRef)
	if e.Flags&0x01 != 0 {
		enc.Float32(e.CenterOffset.X)
		enc.Float32(e.CenterOffset.Y)
		enc.Float32(e.CenterOffset.Z)
	}
	if e.Flags&0x02 != 0 {
		enc.Float32(e.Radius)
	}
	if e.Flags&0x04 != 0 {
		enc.Uint32(e.CurrentFrame)
	}
	if e.Flags&0x08 != 0 {
		enc.Uint32(e.Sleep)
	}
	if e.Flags&0x10 != 0 {
		for _, spriteFragment := range e.SpriteFragments {
			enc.Uint32(spriteFragment)
		}
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSprite4DDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	frameCount := dec.Uint32()
	e.PolyRef = dec.Int32()
	if e.Flags&0x01 != 0 {
		e.CenterOffset.X = dec.Float32()
		e.CenterOffset.Y = dec.Float32()
		e.CenterOffset.Z = dec.Float32()
	}
	if e.Flags&0x02 != 0 {
		e.Radius = dec.Float32()
	}
	if e.Flags&0x04 != 0 {
		e.CurrentFrame = dec.Uint32()
	}
	if e.Flags&0x08 != 0 {
		e.Sleep = dec.Uint32()
	}
	if e.Flags&0x10 != 0 {
		for i := uint32(0); i < frameCount; i++ {
			e.SpriteFragments = append(e.SpriteFragments, dec.Uint32())
		}
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
