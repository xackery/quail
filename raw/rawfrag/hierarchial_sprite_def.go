package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

////////

// WldFragHierarchialSpriteDef is HierarchialSpriteDef in libeq, SkeletonTrackSet in openzone, HIERARCHIALSPRITE in wld, SkeletonHierarchy in lantern
type WldFragHierarchialSpriteDef struct {
	NameRef            int32                  `yaml:"name_ref"`
	Flags              uint32                 `yaml:"flags"`
	CollisionVolumeRef uint32                 `yaml:"collision_volume_ref"`
	CenterOffset       [3]uint32              `yaml:"center_offset"`
	BoundingRadius     float32                `yaml:"bounding_radius"`
	Bones              []WldFragSkeletonEntry `yaml:"bones"`
	Skins              []uint32               `yaml:"skins"`
	SkinLinks          []uint32               `yaml:"skin_links"`
}

type WldFragSkeletonEntry struct {
	NameRef                   int32    `yaml:"name_ref"`
	Flags                     uint32   `yaml:"flags"`
	TrackRef                  uint32   `yaml:"track_ref"`
	MeshOrSpriteOrParticleRef uint32   `yaml:"mesh_or_sprite_or_particle_ref"`
	SubBones                  []uint32 `yaml:"sub_bones"`
}

func (e *WldFragHierarchialSpriteDef) FragCode() int {
	return FragCodeHierarchialSpriteDef
}

func (e *WldFragHierarchialSpriteDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Bones)))
	enc.Uint32(e.CollisionVolumeRef)
	if e.Flags&0x1 != 0 {
		enc.Uint32(e.CenterOffset[0])
		enc.Uint32(e.CenterOffset[1])
		enc.Uint32(e.CenterOffset[2])
	}

	if e.Flags&0x2 != 0 {
		enc.Float32(e.BoundingRadius)
	}

	for _, bone := range e.Bones {
		enc.Int32(bone.NameRef)
		enc.Uint32(bone.Flags)
		enc.Uint32(bone.TrackRef)
		enc.Uint32(bone.MeshOrSpriteOrParticleRef)
		enc.Uint32(uint32(len(bone.SubBones)))
		for _, subCount := range bone.SubBones {
			enc.Uint32(subCount)
		}
	}

	if e.Flags&0x200 != 0 {
		enc.Uint32(uint32(len(e.Skins)))
		for _, skin := range e.Skins {
			enc.Uint32(skin)
		}
		for _, skinLink := range e.SkinLinks {
			enc.Uint32(skinLink)
		}
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragHierarchialSpriteDef) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	boneCount := dec.Uint32()
	e.CollisionVolumeRef = dec.Uint32()
	if e.Flags&0x1 != 0 {
		e.CenterOffset = [3]uint32{dec.Uint32(), dec.Uint32(), dec.Uint32()}
	}

	if e.Flags&0x2 != 0 {
		e.BoundingRadius = dec.Float32()
	}

	for i := 0; i < int(boneCount); i++ {
		bone := WldFragSkeletonEntry{}
		bone.NameRef = dec.Int32()
		bone.Flags = dec.Uint32()
		bone.TrackRef = dec.Uint32()
		bone.MeshOrSpriteOrParticleRef = dec.Uint32()
		subCount := dec.Uint32()
		for j := 0; j < int(subCount); j++ {
			bone.SubBones = append(bone.SubBones, dec.Uint32())
		}
		e.Bones = append(e.Bones, bone)
	}

	if e.Flags&0x200 != 0 {
		skinCount := dec.Uint32()
		for i := 0; i < int(skinCount); i++ {
			e.Skins = append(e.Skins, dec.Uint32())
		}
		for i := 0; i < int(skinCount); i++ {
			e.SkinLinks = append(e.SkinLinks, dec.Uint32())
		}
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
