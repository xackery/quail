package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

const (
	ActorFlagHasCurrentAction  = 0x0001
	ActorFlagHasLocation       = 0x0002
	ActorFlagHasBoundingRadius = 0x0004
	ActorFlagHasScaleFactor    = 0x0008
	ActorFlagHasSound          = 0x0010
	ActorFlagActive            = 0x0020
	ActorFlagActiveGeometry    = 0x0040
	ActorFlagSpriteVolumeOnly  = 0x0080
	ActorFlagHaveDMRGBTrack    = 0x0100
	ActorFlagUsesBoundingBox   = 0x0200
)

// WldFragActorDef is ActorDef in libeq, Static in openzone, ACTORDEF in wld
type WldFragActorDef struct {
	nameRef         int32
	Flags           uint32
	CallbackNameRef int32
	BoundsRef       int32 // ref to sphere, spherelist or polyhedron
	CurrentAction   uint32
	Location        [6]float32
	Unk1            uint32
	Actions         []WldFragModelAction
	SpriteRefs      []uint32
	UserData        string
}

type WldFragModelAction struct {
	Unk1 uint32
	Lods []float32
}

func (e *WldFragActorDef) FragCode() int {
	return FragCodeActorDef
}

func (e *WldFragActorDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	userData := helper.WriteStringHash(e.UserData)
	paddingSize := (4 - (len(userData) % 4)) % 4

	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)

	enc.Int32(e.CallbackNameRef)
	enc.Uint32(uint32(len(e.Actions)))
	enc.Uint32(uint32(len(e.SpriteRefs)))
	enc.Int32(e.BoundsRef)
	if helper.HasFlag(e.Flags, ActorFlagHasCurrentAction) {
		enc.Uint32(e.CurrentAction)
	}
	if helper.HasFlag(e.Flags, ActorFlagHasLocation) {
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
	for _, fragmentRef := range e.SpriteRefs {
		enc.Uint32(fragmentRef)
	}
	enc.Uint32(uint32(len(userData)))
	if len(e.UserData) > 0 {
		enc.Bytes(userData)
		enc.Bytes(make([]byte, paddingSize))
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragActorDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.CallbackNameRef = dec.Int32()
	actionCount := dec.Uint32()
	spriteCount := dec.Uint32()
	e.BoundsRef = dec.Int32()
	if helper.HasFlag(e.Flags, ActorFlagHasCurrentAction) {
		e.CurrentAction = dec.Uint32()
	}
	if helper.HasFlag(e.Flags, ActorFlagHasLocation) {
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
	for i := uint32(0); i < spriteCount; i++ {
		e.SpriteRefs = append(e.SpriteRefs, dec.Uint32())
	}
	userDataSize := dec.Uint32()
	if userDataSize > 0 {
		e.UserData = helper.ReadStringHash([]byte(dec.StringFixed(int(userDataSize))))
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragActorDef) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragActorDef) SetNameRef(nameRef int32) {
	e.nameRef = nameRef
}
