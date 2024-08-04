package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragActorDef is ActorDef in libeq, Static in openzone, ACTORDEF in wld
type WldFragActorDef struct {
	NameRef         int32 `yaml:"name_ref"`
	Flags           uint32
	CallbackNameRef int32  `yaml:"callback_name_ref"`
	BoundsRef       int32  // ref to sphere, spherelist or polyhedron
	CurrentAction   uint32 `yaml:"current_action"`
	Location        [6]float32
	Unk1            uint32               `yaml:"unk1"`
	Actions         []WldFragModelAction `yaml:"actions"`
	FragmentRefs    []uint32             `yaml:"fragment_refs"`
	Unk2            uint32               `yaml:"unk2"`
}

type WldFragModelAction struct {
	Unk1 uint32    `yaml:"unk1"`
	Lods []float32 `yaml:"lods"`
}

func (e *WldFragActorDef) FragCode() int {
	return FragCodeActorDef
}

func (e *WldFragActorDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)

	enc.Int32(e.CallbackNameRef)
	enc.Uint32(uint32(len(e.Actions)))
	enc.Uint32(uint32(len(e.FragmentRefs)))
	enc.Int32(e.BoundsRef)
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
	for _, action := range e.Actions {
		enc.Uint32(uint32(len(action.Lods)))
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
	actionCount := dec.Uint32()
	fragmentRefCount := dec.Uint32()
	e.BoundsRef = dec.Int32()
	if e.Flags&0x01 == 0x01 {
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
	for i := uint32(0); i < actionCount; i++ {
		var action WldFragModelAction
		lodCount := dec.Uint32()
		action.Unk1 = dec.Uint32()
		for j := uint32(0); j < lodCount; j++ {
			action.Lods = append(action.Lods, dec.Float32())
		}
		e.Actions = append(e.Actions, action)
	}
	for i := uint32(0); i < fragmentRefCount; i++ {
		e.FragmentRefs = append(e.FragmentRefs, dec.Uint32())
	}
	e.Unk2 = dec.Uint32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
