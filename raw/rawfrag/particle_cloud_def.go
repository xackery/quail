package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

const (
	ParticleCloudFlagHasSpawnBox  = 0x01
	ParticleCloudFlagHasBox       = 0x02
	ParticleCloudFlagHasSpriteDef = 0x04
)

// WldFragParticleCloudDef is ParticleCloudDef in libeq, empty in openzone, empty in wld, WldFragParticleCloudDef in lantern
type WldFragParticleCloudDef struct {
	nameRef                 int32
	Flags                   uint32
	ParticleType            uint32
	SpawnType               uint32 // 0x01 sphere, 0x02 plane, 0x03 stream, 0x04 none
	PCloudFlags             uint32 //Flag 1, High Opacity, Flag 3, Follows Item
	Size                    uint32
	GravityMultiplier       float32
	Gravity                 [3]float32
	Duration                uint32
	SpawnRadius             float32 // sphere radius
	SpawnAngle              float32 // cone angle
	Lifespan                uint32
	SpawnVelocityMultiplier float32
	SpawnVelocity           [3]float32
	SpawnRate               uint32
	SpawnScale              float32
	Tint                    [4]uint8
	SpawnBoxMin             [3]float32
	SpawnBoxMax             [3]float32
	BoxMin                  [3]float32
	BoxMax                  [3]float32
	BlitSpriteRef           uint32
}

func (e *WldFragParticleCloudDef) FragCode() int {
	return FragCodeParticleCloudDef
}

func (e *WldFragParticleCloudDef) Write(w io.Writer, isNewWorld bool) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.nameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.ParticleType)
	enc.Uint32(e.SpawnType)
	enc.Uint32(e.PCloudFlags)
	enc.Uint32(e.Size)
	enc.Float32(e.GravityMultiplier)
	enc.Float32(e.Gravity[0])
	enc.Float32(e.Gravity[1])
	enc.Float32(e.Gravity[2])
	enc.Uint32(e.Duration)
	enc.Float32(e.SpawnRadius)
	enc.Float32(e.SpawnAngle)
	enc.Uint32(e.Lifespan)
	enc.Float32(e.SpawnVelocityMultiplier)
	enc.Float32(e.SpawnVelocity[0])
	enc.Float32(e.SpawnVelocity[1])
	enc.Float32(e.SpawnVelocity[2])
	enc.Uint32(e.SpawnRate)
	enc.Float32(e.SpawnScale)
	enc.Uint8(e.Tint[0])
	enc.Uint8(e.Tint[1])
	enc.Uint8(e.Tint[2])
	enc.Uint8(e.Tint[3])
	if e.Flags&ParticleCloudFlagHasSpawnBox != 0 {
		enc.Float32(e.SpawnBoxMin[0])
		enc.Float32(e.SpawnBoxMin[1])
		enc.Float32(e.SpawnBoxMin[2])
		enc.Float32(e.SpawnBoxMax[0])
		enc.Float32(e.SpawnBoxMax[1])
		enc.Float32(e.SpawnBoxMax[2])
	}
	if e.Flags&ParticleCloudFlagHasBox != 0 {
		enc.Float32(e.BoxMin[0])
		enc.Float32(e.BoxMin[1])
		enc.Float32(e.BoxMin[2])
		enc.Float32(e.BoxMax[0])
		enc.Float32(e.BoxMax[1])
		enc.Float32(e.BoxMax[2])
	}

	if e.Flags&ParticleCloudFlagHasSpriteDef != 0 {
		enc.Uint32(e.BlitSpriteRef)
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragParticleCloudDef) Read(r io.ReadSeeker, isNewWorld bool) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.nameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.ParticleType = dec.Uint32()
	e.SpawnType = dec.Uint32()
	e.PCloudFlags = dec.Uint32()
	e.Size = dec.Uint32()
	e.GravityMultiplier = dec.Float32()
	e.Gravity = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	e.Duration = dec.Uint32()
	e.SpawnRadius = dec.Float32()
	e.SpawnAngle = dec.Float32()
	e.Lifespan = dec.Uint32()
	e.SpawnVelocityMultiplier = dec.Float32()
	e.SpawnVelocity = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	e.SpawnRate = dec.Uint32()
	e.SpawnScale = dec.Float32()
	e.Tint = [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()}

	if e.Flags&ParticleCloudFlagHasSpawnBox != 0 {
		e.SpawnBoxMin = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
		e.SpawnBoxMax = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	}
	if e.Flags&ParticleCloudFlagHasBox != 0 {
		e.BoxMin = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
		e.BoxMax = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	}
	if e.Flags&ParticleCloudFlagHasSpriteDef != 0 {
		e.BlitSpriteRef = dec.Uint32()
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (e *WldFragParticleCloudDef) NameRef() int32 {
	return e.nameRef
}

func (e *WldFragParticleCloudDef) SetNameRef(id int32) {
	e.nameRef = id
}
