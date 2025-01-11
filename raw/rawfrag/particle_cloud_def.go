package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// WldFragParticleCloudDef is ParticleCloudDef in libeq, empty in openzone, empty in wld, WldFragParticleCloudDef in lantern
type WldFragParticleCloudDef struct {
	parents               []common.TreeLinker
	children              []common.TreeLinker
	fragID                int
	tag                   string
	NameRef               int32
	SettingOne            uint32
	SettingTwo            uint32
	ParticleMovement      uint32 // 0x01 sphere, 0x02 plane, 0x03 stream, 0x04 none
	Flags                 uint32 //Flag 1, High Opacity, Flag 3, Follows Item
	SimultaneousParticles uint32
	Unk6                  uint32
	Unk7                  uint32
	Unk8                  uint32
	Unk9                  uint32
	Unk10                 uint32
	SpawnRadius           float32 // sphere radius
	SpawnAngle            float32 // cone angle
	SpawnLifespan         uint32
	SpawnVelocity         float32
	SpawnNormal           [3]float32
	SpawnRate             uint32
	SpawnScale            float32
	Color                 [4]uint8
	BlitSpriteRef         uint32
}

func (e *WldFragParticleCloudDef) FragCode() int {
	return FragCodeParticleCloudDef
}

func (e *WldFragParticleCloudDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.SettingOne)
	enc.Uint32(e.SettingTwo)
	enc.Uint32(e.ParticleMovement)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SimultaneousParticles)
	enc.Uint32(e.Unk6)
	enc.Uint32(e.Unk7)
	enc.Uint32(e.Unk8)
	enc.Uint32(e.Unk9)
	enc.Uint32(e.Unk10)
	enc.Float32(e.SpawnRadius)
	enc.Float32(e.SpawnAngle)
	enc.Uint32(e.SpawnLifespan)
	enc.Float32(e.SpawnVelocity)
	enc.Float32(e.SpawnNormal[0])
	enc.Float32(e.SpawnNormal[1])
	enc.Float32(e.SpawnNormal[2])

	enc.Uint32(e.SpawnRate)
	enc.Float32(e.SpawnScale)

	enc.Uint8(e.Color[0])
	enc.Uint8(e.Color[1])
	enc.Uint8(e.Color[2])
	enc.Uint8(e.Color[3])

	enc.Uint32(e.BlitSpriteRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragParticleCloudDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.SettingOne = dec.Uint32()
	e.SettingTwo = dec.Uint32()
	e.ParticleMovement = dec.Uint32()
	e.Flags = dec.Uint32()
	e.SimultaneousParticles = dec.Uint32()
	e.Unk6 = dec.Uint32()
	e.Unk7 = dec.Uint32()
	e.Unk8 = dec.Uint32()
	e.Unk9 = dec.Uint32()
	e.Unk10 = dec.Uint32()
	e.SpawnRadius = dec.Float32()
	e.SpawnAngle = dec.Float32()
	e.SpawnLifespan = dec.Uint32()
	e.SpawnVelocity = dec.Float32()
	e.SpawnNormal = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	e.SpawnRate = dec.Uint32()
	e.SpawnScale = dec.Float32()
	e.Color = [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()}
	e.BlitSpriteRef = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragParticleCloudDef) Parents() []common.TreeLinker {
	return e.parents
}

func (e *WldFragParticleCloudDef) AddParent(parent common.TreeLinker) {
	e.parents = append(e.parents, parent)
}

func (e *WldFragParticleCloudDef) Tag() string {
	return e.tag
}

func (e *WldFragParticleCloudDef) SetFragID(id int) {
	e.fragID = id
}

func (e *WldFragParticleCloudDef) FragID() int {
	return e.fragID
}

func (e *WldFragParticleCloudDef) Children() []common.TreeLinker {
	return nil
}

func (e *WldFragParticleCloudDef) FragType() string {
	return "PACD"
}

func (e *WldFragParticleCloudDef) AddChild(child common.TreeLinker) {
	e.children = append(e.children, child)
}
