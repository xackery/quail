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

// WldFragHierarchialSprite is HierarchialSprite in libeq, SkeletonTrackSetReference in openzone, HIERARCHIALSPRITE (ref) in wld, SkeletonHierarchyReference in lantern
type WldFragHierarchialSprite struct {
	NameRef              int16  `yaml:"name_ref"`
	HierarchialSpriteRef int16  `yaml:"hierarchial_sprite_ref"`
	Flags                uint32 `yaml:"flags"`
}

func (e *WldFragHierarchialSprite) FragCode() int {
	return FragCodeHierarchialSprite
}

func (e *WldFragHierarchialSprite) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int16(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Int16(e.HierarchialSpriteRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragHierarchialSprite) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int16()
	e.Flags = dec.Uint32()
	e.HierarchialSpriteRef = dec.Int16()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// WldFragTrackDef is TrackDef in libeq, Mob Skeleton Piece WldFragTrackDef in openzone, TRACKDEFINITION in wld, TrackDefFragment in lantern
type WldFragTrackDef struct {
	NameRef        int32                       `yaml:"name_ref"`
	Flags          uint32                      `yaml:"flags"`
	BoneTransforms []WldFragTrackBoneTransform `yaml:"skeleton_transforms"`
}

type WldFragTrackBoneTransform struct {
	RotateDenominator int16 `yaml:"rotate_denominator"`
	RotateX           int16 `yaml:"rotate_x"`
	RotateY           int16 `yaml:"rotate_y"`
	RotateZ           int16 `yaml:"rotate_z"`
	ShiftDenominator  int16 `yaml:"shift_denominator"`
	ShiftX            int16 `yaml:"shift_x"`
	ShiftY            int16 `yaml:"shift_y"`
	ShiftZ            int16 `yaml:"shift_z"`
}

func (e *WldFragTrackDef) FragCode() int {
	return FragCodeTrackDef
}

func (e *WldFragTrackDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.BoneTransforms)))
	for _, ft := range e.BoneTransforms {
		if e.Flags&0x08 == 0x08 {
			enc.Int16(ft.ShiftDenominator)
			enc.Int16(ft.ShiftX)
			enc.Int16(ft.ShiftY)
			enc.Int16(ft.ShiftZ)
			enc.Int16(ft.RotateX)
			enc.Int16(ft.RotateY)
			enc.Int16(ft.RotateZ)
			enc.Int16(ft.RotateDenominator)
			continue
		}
		enc.Int8(int8(ft.RotateDenominator))
		enc.Int8(int8(ft.RotateX))
		enc.Int8(int8(ft.RotateY))
		enc.Int8(int8(ft.RotateZ))
		enc.Int8(int8(ft.ShiftDenominator))
		enc.Int8(int8(ft.ShiftX))
		enc.Int8(int8(ft.ShiftY))
		enc.Int8(int8(ft.ShiftZ))

	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragTrackDef) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	boneCount := dec.Uint32()
	for i := 0; i < int(boneCount); i++ {
		ft := WldFragTrackBoneTransform{}
		if e.Flags&0x08 == 0x08 {
			ft.ShiftDenominator = dec.Int16()
			ft.ShiftX = dec.Int16()
			ft.ShiftY = dec.Int16()
			ft.ShiftZ = dec.Int16()
			ft.RotateX = dec.Int16()
			ft.RotateY = dec.Int16()
			ft.RotateZ = dec.Int16()
			ft.RotateDenominator = dec.Int16()
			e.BoneTransforms = append(e.BoneTransforms, ft)
			continue
		}
		ft.RotateDenominator = int16(dec.Int8())
		ft.RotateX = int16(dec.Int8())
		ft.RotateY = int16(dec.Int8())
		ft.RotateZ = int16(dec.Int8())
		ft.ShiftDenominator = int16(dec.Int8())
		ft.ShiftX = int16(dec.Int8())
		ft.ShiftY = int16(dec.Int8())
		ft.ShiftZ = int16(dec.Int8())
		e.BoneTransforms = append(e.BoneTransforms, ft)

	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// WldFragTrack is a bone in a skeleton. It is Track in libeq, Mob Skeleton Piece Track Reference in openzone, TRACKINSTANCE in wld, TrackDefFragment in lantern
type WldFragTrack struct {
	NameRef int32  `yaml:"name_ref"`
	Track   int32  `yaml:"track_ref"`
	Flags   uint32 `yaml:"flags"`
	Sleep   uint32 `yaml:"sleep"` // if 0x01 is set, this is the number of milliseconds to sleep before starting the animation
}

func (e *WldFragTrack) FragCode() int {
	return FragCodeTrack
}

func (e *WldFragTrack) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.Track)
	enc.Uint32(e.Flags)
	if e.Flags&0x01 == 0x01 {
		enc.Uint32(e.Sleep)
	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragTrack) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Track = dec.Int32()
	e.Flags = dec.Uint32()
	if e.Flags&0x01 == 0x01 {
		e.Sleep = dec.Uint32()
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// WldFragDMTrack is DmTrackDef in libeq, empty in openzone, empty in wld
type WldFragDMTrackDef struct {
}

func (e *WldFragDMTrackDef) FragCode() int {
	return FragCodeDMTrackDef
}

func (e *WldFragDMTrackDef) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDMTrackDef) Read(r io.ReadSeeker) error {
	return nil
}

// WldFragDMTrack is DmTrack in libeq, Mesh Animated Vertices Reference in openzone, empty in wld, MeshAnimatedVerticesReference in lantern
type WldFragDMTrack struct {
}

func (e *WldFragDMTrack) FragCode() int {
	return FragCodeDMTrack
}

func (e *WldFragDMTrack) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDMTrack) Read(r io.ReadSeeker) error {
	return nil
}

// WldFragDmRGBTrackDef is a list of colors, one per vertex, for baked lighting. It is DmRGBTrackDef in libeq, Vertex Color in openzone, empty in wld, VertexColors in lantern
type WldFragDmRGBTrackDef struct {
}

func (e *WldFragDmRGBTrackDef) FragCode() int {
	return FragCodeDmRGBTrackDef
}

func (e *WldFragDmRGBTrackDef) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDmRGBTrackDef) Read(r io.ReadSeeker) error {
	return nil
}

// WldFragDmRGBTrack is DmRGBTrack in libeq, Vertex Color Reference in openzone, empty in wld, VertexColorsReference in lantern
type WldFragDmRGBTrack struct {
}

func (e *WldFragDmRGBTrack) FragCode() int {
	return FragCodeDmRGBTrack
}

func (e *WldFragDmRGBTrack) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDmRGBTrack) Read(r io.ReadSeeker) error {
	return nil
}
