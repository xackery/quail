package vwld

// VWld is a struct representing a VWld file
type VWld struct {
	FileName           string
	Version            uint32
	Bitmaps            []*Bitmap
	Sprites            []*Sprite
	SpriteInstances    []*SpriteInstance
	Particles          []*Particle
	ParticleInstances  []*ParticleInstance
	Materials          []*Material
	MaterialInstances  []*MaterialInstance
	Meshes             []*Mesh
	Animations         []*Animation
	AnimationInstances []*AnimationInstance
	MeshInstances      []*MeshInstance
}

// Bitmap is a struct representing a material
type Bitmap struct {
	fragID          uint32
	Name            string
	Textures        []string
	SimpleSpriteDef Sprite
	SimpleSprite    SpriteInstance
}

func (wld *VWld) bitmapByFragID(fragID uint32) *Bitmap {
	for _, bitmap := range wld.Bitmaps {
		if bitmap.fragID == fragID {
			return bitmap
		}
	}
	return nil
}

type Sprite struct {
	fragID       uint32
	Name         string
	Flags        uint32
	CurrentFrame int32
	Sleep        uint32
	Bitmaps      []string
}

func (wld *VWld) spriteByFragID(fragID uint32) *Sprite {
	for _, sprite := range wld.Sprites {
		if sprite.fragID == fragID {
			return sprite
		}
	}
	return nil
}

type SpriteInstance struct {
	fragID uint32
	Name   string
	Flags  uint32
	Sprite string
}

func (wld *VWld) spriteInstanceByFragID(fragID uint32) *SpriteInstance {
	for _, spriteInstance := range wld.SpriteInstances {
		if spriteInstance.fragID == fragID {
			return spriteInstance
		}
	}
	return nil
}

// Particle is also known as BlitSpriteDef
type Particle struct {
	fragID           uint32
	Name             string
	Flags            uint32
	SpriteName       string
	Unknown          int32
	ParticleCloudDef ParticleInstance
}

func (wld *VWld) particleByFragID(fragID uint32) *Particle {
	for _, particle := range wld.Particles {
		if particle.fragID == fragID {
			return particle
		}
	}
	return nil
}

type ParticleInstance struct {
	fragID                uint32
	Name                  string  `yaml:"name"`
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
	Particle              string  `yaml:"particle"`
}

func (wld *VWld) particleInstanceByFragID(fragID uint32) *ParticleInstance {
	for _, particleInstance := range wld.ParticleInstances {
		if particleInstance.fragID == fragID {
			return particleInstance
		}
	}
	return nil
}

// Material is a struct representing a material
type Material struct {
	fragID        uint32
	Name          string
	Flags         uint32    `yaml:"flags"`
	RenderMethod  uint32    `yaml:"render_method"`
	RGBPen        uint32    `yaml:"rgb_pen"`
	Brightness    float32   `yaml:"brightness"`
	ScaledAmbient float32   `yaml:"scaled_ambient"`
	Texture       string    `yaml:"texture"`
	Pairs         [2]uint32 `yaml:"pairs"`
	Palette       MaterialInstance
}

func (wld *VWld) materialByFragID(fragID uint32) *Material {
	for _, material := range wld.Materials {
		if material.fragID == fragID {
			return material
		}
	}
	return nil
}

// MaterialInstance is a struct representing a material palette
type MaterialInstance struct {
	fragID    uint32
	Name      string
	Flags     uint32
	Materials []string
}

func (wld *VWld) materialInstanceByFragID(fragID uint32) *MaterialInstance {
	for _, materialPalette := range wld.MaterialInstances {
		if materialPalette.fragID == fragID {
			return materialPalette
		}
	}
	return nil
}

type Mesh struct {
	fragID            uint32
	Name              string
	Flags             uint32     `yaml:"flags"`
	MaterialInstance  string     `yaml:"material_instance"`
	AnimationInstance string     `yaml:"animation_instance"`
	Fragment3Ref      int32      `yaml:"fragment_3_ref"`
	Fragment4Ref      int32      `yaml:"fragment_4_ref"` // unknown, usually ref to first texture
	Center            Vector3    `yaml:"center"`
	Params2           UIndex3    `yaml:"params_2"`
	MaxDistance       float32    `yaml:"max_distance"`
	Min               Vector3    `yaml:"min"`
	Max               Vector3    `yaml:"max"`
	RawScale          uint16     `yaml:"raw_scale"`
	MeshopCount       uint16     `yaml:"meshop_count"`
	Scale             float32    `yaml:"scale"`
	Vertices          [][3]int16 `yaml:"vertices"`
	UVs               [][2]int16 `yaml:"uvs"`
	Normals           [][3]int8  `yaml:"normals"`
	Colors            []RGBA     `yaml:"colors"`
	Triangles         []Triangle `yaml:"triangles"`
}

func (wld *VWld) meshByFragID(fragID uint32) *Mesh {
	for _, mesh := range wld.Meshes {
		if mesh.fragID == fragID {
			return mesh
		}
	}
	return nil
}

type MeshInstance struct {
	fragID uint32
	Name   string
	Mesh   string
	Params uint32
}

func (wld *VWld) meshInstanceByFragID(fragID uint32) *MeshInstance {
	for _, meshInstance := range wld.MeshInstances {
		if meshInstance.fragID == fragID {
			return meshInstance
		}
	}
	return nil
}

type Animation struct {
	fragID     uint32
	Name       string
	Flags      uint32
	Transforms []*AnimationTransform
}

type AnimationTransform struct {
	RotateDenominator int16
	RotateX           int16
	RotateY           int16
	RotateZ           int16
	ShiftX            int16
	ShiftY            int16
	ShiftZ            int16
	ShiftDenominator  int16
}

func (wld *VWld) animationByFragID(fragID uint32) *Animation {
	for _, animation := range wld.Animations {
		if animation.fragID == fragID {
			return animation
		}
	}
	return nil
}

type AnimationInstance struct {
	fragID    uint32
	Name      string
	Animation string
	Flags     uint32
	Sleep     uint32
}

func (wld *VWld) animationInstanceByFragID(fragID uint32) *AnimationInstance {
	for _, animationInstance := range wld.AnimationInstances {
		if animationInstance.fragID == fragID {
			return animationInstance
		}
	}
	return nil
}
