package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/model"
)

// WldFragParticleCloudDef is ParticleCloudDef in libeq, empty in openzone, empty in wld, WldFragParticleCloudDef in lantern
type WldFragParticleCloudDef struct {
	NameRef               int32      `yaml:"name_ref"`
	Unk1                  uint32     `yaml:"unk1"`
	Unk2                  uint32     `yaml:"unk2"`
	ParticleMovement      uint32     `yaml:"particle_movement"` // 0x01 sphere, 0x02 plane, 0x03 stream, 0x04 none
	Flags                 uint32     //Flag 1, High Opacity, Flag 3, Follows Item
	SimultaneousParticles uint32     `yaml:"simultaneous_particles"`
	Unk6                  uint32     `yaml:"unk6"`
	Unk7                  uint32     `yaml:"unk7"`
	Unk8                  uint32     `yaml:"unk8"`
	Unk9                  uint32     `yaml:"unk9"`
	Unk10                 uint32     `yaml:"unk10"`
	SpawnRadius           float32    `yaml:"spawn_radius"` // sphere radius
	SpawnAngle            float32    `yaml:"spawn_angle"`  // cone angle
	SpawnLifespan         uint32     `yaml:"spawn_lifespan"`
	SpawnVelocity         float32    `yaml:"spawn_velocity"`
	SpawnNormalZ          float32    `yaml:"spawn_normal_z"`
	SpawnNormalX          float32    `yaml:"spawn_normal_x"`
	SpawnNormalY          float32    `yaml:"spawn_normal_y"`
	SpawnRate             uint32     `yaml:"spawn_rate"`
	SpawnScale            float32    `yaml:"spawn_scale"`
	Color                 model.RGBA `yaml:"color"`
	ParticleRef           uint32     `yaml:"particle_ref"`
}

func (e *WldFragParticleCloudDef) FragCode() int {
	return FragCodeParticleCloudDef
}

func (e *WldFragParticleCloudDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Unk1)
	enc.Uint32(e.Unk2)
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
	enc.Float32(e.SpawnNormalZ)
	enc.Float32(e.SpawnNormalX)
	enc.Float32(e.SpawnNormalY)
	enc.Uint32(e.SpawnRate)
	enc.Float32(e.SpawnScale)
	enc.Uint8(e.Color.R)
	enc.Uint8(e.Color.G)
	enc.Uint8(e.Color.B)
	enc.Uint8(e.Color.A)

	enc.Uint32(e.ParticleRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragParticleCloudDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Unk1 = dec.Uint32()
	e.Unk2 = dec.Uint32()
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
	e.SpawnNormalZ = dec.Float32()
	e.SpawnNormalX = dec.Float32()
	e.SpawnNormalY = dec.Float32()
	e.SpawnRate = dec.Uint32()
	e.SpawnScale = dec.Float32()
	e.Color = model.RGBA{R: dec.Uint8(), G: dec.Uint8(), B: dec.Uint8(), A: dec.Uint8()}
	e.ParticleRef = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
