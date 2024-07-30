package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragActor is Actor in libeq, Object Location in openzone, ACTORINST in wld, ObjectInstance in lantern
type WldFragActor struct {
	NameRef        int32
	ActorDefRef    int32
	Flags          uint32
	SphereRef      uint32
	CurrentAction  uint32
	Location       [6]float32
	Unk1           uint32
	BoundingRadius float32
	ScaleFactor    float32
	SoundNameRef   int32
	DMRGBTrackRef  int32
	UserData       string
}

func (e *WldFragActor) FragCode() int {
	return FragCodeActor
}

func (e *WldFragActor) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.ActorDefRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SphereRef)
	if e.Flags&0x1 == 0x1 {
		enc.Uint32(e.CurrentAction)
	}
	if e.Flags&0x2 == 0x2 {
		enc.Float32(e.Location[0])
		enc.Float32(e.Location[1])
		enc.Float32(e.Location[2])
		enc.Float32(e.Location[3])
		enc.Float32(e.Location[4])
		enc.Float32(e.Location[5])
		enc.Uint32(e.Unk1)
	}
	if e.Flags&0x4 == 0x4 {
		enc.Float32(e.BoundingRadius)
	}
	if e.Flags&0x8 == 0x8 {
		enc.Float32(e.ScaleFactor)
	}
	if e.Flags&0x10 == 0x10 {
		enc.Int32(e.SoundNameRef)
	}
	if e.Flags&0x100 == 0x100 {
		enc.Int32(e.DMRGBTrackRef)
	}

	enc.StringLenPrefixUint32(e.UserData)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragActor) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.ActorDefRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.SphereRef = dec.Uint32()
	if e.Flags&0x1 == 0x1 {
		e.CurrentAction = dec.Uint32()
	}
	if e.Flags&0x2 == 0x2 {
		e.Location[0] = dec.Float32()
		e.Location[1] = dec.Float32()
		e.Location[2] = dec.Float32()
		e.Location[3] = dec.Float32()
		e.Location[4] = dec.Float32()
		e.Location[5] = dec.Float32()
		e.Unk1 = dec.Uint32()
	}
	if e.Flags&0x4 == 0x4 {
		e.BoundingRadius = dec.Float32()
	}
	if e.Flags&0x8 == 0x8 {
		e.ScaleFactor = dec.Float32()
	}
	if e.Flags&0x10 == 0x10 {
		e.SoundNameRef = dec.Int32()
	}
	if e.Flags&0x100 == 0x100 {
		e.DMRGBTrackRef = dec.Int32()
	}
	e.UserData = dec.StringLenPrefixUint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
