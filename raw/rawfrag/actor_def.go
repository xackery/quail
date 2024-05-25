package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragActorDef is ActorDef in libeq, Static in openzone, ACTORDEF in wld
type WldFragActorDef struct {
	NameRef          int32                `yaml:"name_ref"`
	Flags            uint32               `yaml:"flags"`
	CallbackNameRef  int32                `yaml:"callback_name_ref"`
	ActionCount      uint32               `yaml:"action_count"`
	FragmentRefCount uint32               `yaml:"fragment_ref_count"`
	BoundsRef        int32                `yaml:"bounds_ref"`
	CurrentAction    uint32               `yaml:"current_action"`
	Offset           model.Vector3        `yaml:"offset"`
	Rotation         model.Vector3        `yaml:"rotation"`
	Unk1             uint32               `yaml:"unk1"`
	Actions          []WldFragModelAction `yaml:"actions"`
	FragmentRefs     []uint32             `yaml:"fragment_refs"`
	Unk2             uint32               `yaml:"unk2"`
}

type WldFragModelAction struct {
	LodCount uint32    `yaml:"lod_count"`
	Unk1     uint32    `yaml:"unk1"`
	Lods     []float32 `yaml:"lods"`
}

func (e *WldFragActorDef) FragCode() int {
	return FragCodeActorDef
}

func (e *WldFragActorDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)

	enc.Int32(e.CallbackNameRef)
	enc.Uint32(e.ActionCount)
	enc.Uint32(e.FragmentRefCount)
	enc.Int32(e.BoundsRef)
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
	for _, action := range e.Actions {
		enc.Uint32(action.LodCount)
		enc.Uint32(action.Unk1)
		for _, lod := range action.Lods {
			enc.Float32(lod)
		}
	}
	for _, fragmentRef := range e.FragmentRefs {
		enc.Uint32(fragmentRef)
	}
	enc.Uint32(e.Unk2)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragActorDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.CallbackNameRef = dec.Int32()
	e.ActionCount = dec.Uint32()
	e.FragmentRefCount = dec.Uint32()
	e.BoundsRef = dec.Int32()
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
	for i := uint32(0); i < e.ActionCount; i++ {
		var action WldFragModelAction
		action.LodCount = dec.Uint32()
		action.Unk1 = dec.Uint32()
		for j := uint32(0); j < action.LodCount; j++ {
			action.Lods = append(action.Lods, dec.Float32())
		}
		e.Actions = append(e.Actions, action)
	}
	for i := uint32(0); i < e.FragmentRefCount; i++ {
		e.FragmentRefs = append(e.FragmentRefs, dec.Uint32())
	}
	e.Unk2 = dec.Uint32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
