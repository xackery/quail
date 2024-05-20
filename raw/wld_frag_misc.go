package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragDefault is empty in libeq, empty in openzone, DEFAULT?? in wld
type WldFragDefault struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragDefault) FragCode() int {
	return FragCodeDefault
}

func (e *WldFragDefault) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDefault) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	return nil
}

// WldFragGlobalAmbientLightDef is GlobalAmbientLightDef in libeq, WldFragGlobalAmbientLightDef Fragment in openzone, empty in wld, GlobalAmbientLight in lantern
type WldFragGlobalAmbientLightDef struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32
}

func (e *WldFragGlobalAmbientLightDef) FragCode() int {
	return FragCodeGlobalAmbientLightDef
}

// Read writes the fragment to the writer
func (e *WldFragGlobalAmbientLightDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragGlobalAmbientLightDef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragUserData is empty in libeq, empty in openzone, USERDATA in wld
type WldFragUserData struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragUserData) FragCode() int {
	return FragCodeUserData
}

func (e *WldFragUserData) Write(w io.Writer) error {
	return nil
}

func (e *WldFragUserData) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	return nil
}

// WldFragSound is empty in libeq, empty in openzone, SOUNDDEFINITION in wld
type WldFragSound struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	Flags    uint32 `yaml:"flags"`
}

func (e *WldFragSound) FragCode() int {
	return FragCodeSound
}

func (e *WldFragSound) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSound) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragSoundDef is empty in libeq, empty in openzone, SOUNDINSTANCE in wld
type WldFragSoundDef struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	Flags    uint32 `yaml:"flags"`
}

func (e *WldFragSoundDef) FragCode() int {
	return FragCodeSoundDef
}

func (e *WldFragSoundDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSoundDef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragActiveGeoRegion is empty in libeq, empty in openzone, ACTIVEGEOMETRYREGION in wld
type WldFragActiveGeoRegion struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragActiveGeoRegion) FragCode() int {
	return FragCodeActiveGeoRegion
}

func (e *WldFragActiveGeoRegion) Write(w io.Writer) error {
	return nil
}

func (e *WldFragActiveGeoRegion) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	return nil
}

// WldFragSkyRegion is empty in libeq, empty in openzone, SKYREGION in wld
type WldFragSkyRegion struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragSkyRegion) FragCode() int {
	return FragCodeSkyRegion
}

func (e *WldFragSkyRegion) Write(w io.Writer) error {
	return nil
}

func (e *WldFragSkyRegion) Read(r io.ReadSeeker) error {
	return nil
}

// WldFragZone is Zone in libeq, Region Flag in openzone, ZONE in wld, BspRegionType in lantern
type WldFragZone struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragZone) FragCode() int {
	return FragCodeZone
}

func (e *WldFragZone) Write(w io.Writer) error {
	return nil
}

func (e *WldFragZone) Read(r io.ReadSeeker) error {
	return nil
}

// WldFragParticleCloudDef is ParticleCloudDef in libeq, empty in openzone, empty in wld, WldFragParticleCloudDef in lantern
type WldFragParticleCloudDef struct {
	FragName              string  `yaml:"frag_name"`
	NameRef               int32   `yaml:"name_ref"`
	Unk1                  uint32  `yaml:"unk1"`
	Unk2                  uint32  `yaml:"unk2"`
	ParticleMovement      uint32  `yaml:"particle_movement"` // 0x01 sphere, 0x02 plane, 0x03 stream, 0x04 none
	Flags                 uint32  //Flag 1, High Opacity, Flag 3, Follows Item
	SimultaneousParticles uint32  `yaml:"simultaneous_particles"`
	Unk6                  uint32  `yaml:"unk6"`
	Unk7                  uint32  `yaml:"unk7"`
	Unk8                  uint32  `yaml:"unk8"`
	Unk9                  uint32  `yaml:"unk9"`
	Unk10                 uint32  `yaml:"unk10"`
	SpawnRadius           float32 `yaml:"spawn_radius"` // sphere radius
	SpawnAngle            float32 `yaml:"spawn_angle"`  // cone angle
	SpawnLifespan         uint32  `yaml:"spawn_lifespan"`
	SpawnVelocity         float32 `yaml:"spawn_velocity"`
	SpawnNormalZ          float32 `yaml:"spawn_normal_z"`
	SpawnNormalX          float32 `yaml:"spawn_normal_x"`
	SpawnNormalY          float32 `yaml:"spawn_normal_y"`
	SpawnRate             uint32  `yaml:"spawn_rate"`
	SpawnScale            float32 `yaml:"spawn_scale"`
	Color                 RGBA    `yaml:"color"`
	ParticleRef           uint32  `yaml:"particle_ref"`
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
	e.FragName = FragName(e.FragCode())
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
	e.Color = RGBA{R: dec.Uint8(), G: dec.Uint8(), B: dec.Uint8(), A: dec.Uint8()}
	e.ParticleRef = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
