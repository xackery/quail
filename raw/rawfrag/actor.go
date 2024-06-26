package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragActor is Actor in libeq, Object Location in openzone, ACTORINST in wld, ObjectInstance in lantern
type WldFragActor struct {
	NameRef        int32         `yaml:"name_ref"`
	ActorDefRef    int32         `yaml:"actor_def_ref"`
	Flags          uint32        `yaml:"flags"`
	SphereRef      uint32        `yaml:"sphere_ref"`
	CurrentAction  uint32        `yaml:"current_action"`
	Offset         model.Vector3 `yaml:"offset"`
	Rotation       model.Vector3 `yaml:"rotation"`
	Unk1           uint32        `yaml:"unk1"`
	BoundingRadius float32       `yaml:"bounding_radius"`
	Scale          float32       `yaml:"scale"`
	SoundNameRef   int32         `yaml:"sound_name_ref"`
	Unk2           int32         `yaml:"unk2"`
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
		enc.Float32(e.Offset.X)
		enc.Float32(e.Offset.Y)
		enc.Float32(e.Offset.Z)
		enc.Float32(e.Rotation.X)
		enc.Float32(e.Rotation.Y)
		enc.Float32(e.Rotation.Z)
		enc.Uint32(e.Unk1)
	}
	if e.Flags&0x4 == 0x4 {
		enc.Float32(e.BoundingRadius)
	}
	if e.Flags&0x8 == 0x8 {
		enc.Float32(e.Scale)
	}
	if e.Flags&0x10 == 0x10 {
		enc.Int32(e.SoundNameRef)
	}
	enc.Int32(e.Unk2)
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
		e.Offset.X = dec.Float32()
		e.Offset.Y = dec.Float32()
		e.Offset.Z = dec.Float32()
		e.Rotation.X = dec.Float32()
		e.Rotation.Y = dec.Float32()
		e.Rotation.Z = dec.Float32()
		e.Unk1 = dec.Uint32()
	}
	if e.Flags&0x4 == 0x4 {
		e.BoundingRadius = dec.Float32()
	}
	if e.Flags&0x8 == 0x8 {
		e.Scale = dec.Float32()
	}
	if e.Flags&0x10 == 0x10 {
		e.SoundNameRef = dec.Int32()
	}
	e.Unk2 = dec.Int32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
