package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragHierarchicalSpriteDef is HierarchicalSpriteDef in libeq, SkeletonTrackSet in openzone, HIERARCHICALSPRITE in wld, SkeletonHierarchy in lantern
type WldFragHierarchicalSpriteDef struct {
	NameRef                     int32        `yaml:"name_ref"`
	Flags                       uint32       `yaml:"flags"`
	CollisionVolumeRef          uint32       `yaml:"collision_volume_ref"`
	CenterOffset                [3]float32   `yaml:"center_offset"`
	BoundingRadius              float32      `yaml:"bounding_radius"`
	Dags                        []WldFragDag `yaml:"bones"`
	DMSprites                   []uint32
	LinkSkinUpdatesToDagIndexes []uint32 `yaml:"skin_links"`
}

type WldFragDag struct {
	NameRef                   int32    `yaml:"name_ref"`
	Flags                     uint32   `yaml:"flags"`
	TrackRef                  uint32   `yaml:"track_ref"`
	MeshOrSpriteOrParticleRef uint32   `yaml:"mesh_or_sprite_or_particle_ref"`
	SubDags                   []uint32 `yaml:"sub_bones"`
}

func (e *WldFragHierarchicalSpriteDef) FragCode() int {
	return FragCodeHierarchicalSpriteDef
}

func (e *WldFragHierarchicalSpriteDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.Dags)))
	enc.Uint32(e.CollisionVolumeRef)
	if e.Flags&0x1 != 0 {
		enc.Float32(e.CenterOffset[0])
		enc.Float32(e.CenterOffset[1])
		enc.Float32(e.CenterOffset[2])
	}

	if e.Flags&0x2 != 0 {
		enc.Float32(e.BoundingRadius)
	}

	for _, bone := range e.Dags {
		enc.Int32(bone.NameRef)
		enc.Uint32(bone.Flags)
		enc.Uint32(bone.TrackRef)
		enc.Uint32(bone.MeshOrSpriteOrParticleRef)
		enc.Uint32(uint32(len(bone.SubDags)))
		for _, subCount := range bone.SubDags {
			enc.Uint32(subCount)
		}
	}

	if e.Flags&0x200 != 0 {
		enc.Uint32(uint32(len(e.DMSprites)))
		for _, dmSprite := range e.DMSprites {
			enc.Uint32(dmSprite)
		}
		for _, skinLink := range e.LinkSkinUpdatesToDagIndexes {
			enc.Uint32(skinLink)
		}
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragHierarchicalSpriteDef) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	numDags := dec.Uint32()
	e.CollisionVolumeRef = dec.Uint32()
	if e.Flags&0x1 != 0 {
		e.CenterOffset = [3]float32{dec.Float32(), dec.Float32(), dec.Float32()}
	}

	if e.Flags&0x2 != 0 {
		e.BoundingRadius = dec.Float32()
	}

	for i := 0; i < int(numDags); i++ {
		dag := WldFragDag{}
		dag.NameRef = dec.Int32()
		dag.Flags = dec.Uint32()
		dag.TrackRef = dec.Uint32()
		dag.MeshOrSpriteOrParticleRef = dec.Uint32()
		subCount := dec.Uint32()
		for j := 0; j < int(subCount); j++ {
			dag.SubDags = append(dag.SubDags, dec.Uint32())
		}
		e.Dags = append(e.Dags, dag)
	}

	if e.Flags&0x200 != 0 {
		skinCount := dec.Uint32()
		for i := 0; i < int(skinCount); i++ {
			e.DMSprites = append(e.DMSprites, dec.Uint32())
		}
		for i := 0; i < int(skinCount); i++ {
			e.LinkSkinUpdatesToDagIndexes = append(e.LinkSkinUpdatesToDagIndexes, dec.Uint32())
		}
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
