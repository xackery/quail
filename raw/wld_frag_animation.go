package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragSkeletonTrack is HierarchialSpriteDef in libeq, SkeletonTrackSet in openzone, HIERARCHIALSPRITE in wld, SkeletonHierarchy in lantern
type WldFragSkeletonTrack struct {
	FragName           string                 `yaml:"frag_name"`
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
	NameRef         int32    `yaml:"name_ref"`
	Flags           uint32   `yaml:"flags"`
	TrackRef        uint32   `yaml:"track_ref"`
	MeshOrSpriteRef uint32   `yaml:"mesh_or_sprite_ref"`
	SubBones        []uint32 `yaml:"sub_bones"`
}

func (e *WldFragSkeletonTrack) FragCode() int {
	return 0x10
}

func (e *WldFragSkeletonTrack) Write(w io.Writer) error {
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
		enc.Uint32(bone.MeshOrSpriteRef)
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

func (e *WldFragSkeletonTrack) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())

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
		bone.MeshOrSpriteRef = dec.Uint32()
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

// WldFragSkeletonTrackRef is HierarchialSprite in libeq, SkeletonTrackSetReference in openzone, HIERARCHIALSPRITE (ref) in wld, SkeletonHierarchyReference in lantern
type WldFragSkeletonTrackRef struct {
	FragName         string `yaml:"frag_name"`
	NameRef          int16  `yaml:"name_ref"`
	SkeletonTrackRef int16  `yaml:"skeleton_track_ref"`
	Flags            uint32 `yaml:"flags"`
}

func (e *WldFragSkeletonTrackRef) FragCode() int {
	return 0x11
}

func (e *WldFragSkeletonTrackRef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int16(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Int16(e.SkeletonTrackRef)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSkeletonTrackRef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int16()
	e.Flags = dec.Uint32()
	e.SkeletonTrackRef = dec.Int16()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// WldFragTrack is TrackDef in libeq, Mob Skeleton Piece WldFragTrack in openzone, TRACKDEFINITION in wld, TrackDefFragment in lantern
type WldFragTrack struct {
	FragName      string                      `yaml:"frag_name"`
	NameRef       int32                       `yaml:"name_ref"`
	Flags         uint32                      `yaml:"flags"`
	BoneTransform []WldFragTrackBoneTransform `yaml:"skeleton_transforms"`
}

type WldFragTrackBoneTransform struct {
	TranslationDenominator int16 `yaml:"translation_denominator"`
	TranslationX           int16 `yaml:"translation_x"`
	TranslationY           int16 `yaml:"translation_y"`
	TranslationZ           int16 `yaml:"translation_z"`
	RotationDenominator    int16 `yaml:"rotation_denominator"`
	RotationX              int16 `yaml:"rotation_x"`
	RotationY              int16 `yaml:"rotation_y"`
	RotationZ              int16 `yaml:"rotation_z"`
}

func (e *WldFragTrack) FragCode() int {
	return 0x12
}

func (e *WldFragTrack) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.BoneTransform)))
	for _, ft := range e.BoneTransform {
		if e.Flags&0x08 == 0x08 {
			enc.Int16(ft.RotationDenominator)
			enc.Int16(ft.RotationX)
			enc.Int16(ft.RotationY)
			enc.Int16(ft.RotationZ)
			enc.Int16(ft.TranslationX)
			enc.Int16(ft.TranslationY)
			enc.Int16(ft.TranslationZ)
			enc.Int16(ft.TranslationDenominator)
			continue
		}
		enc.Int8(int8(ft.TranslationDenominator))
		enc.Int8(int8(ft.TranslationX))
		enc.Int8(int8(ft.TranslationY))
		enc.Int8(int8(ft.TranslationZ))
		enc.Int8(int8(ft.RotationDenominator))
		enc.Int8(int8(ft.RotationX))
		enc.Int8(int8(ft.RotationY))
		enc.Int8(int8(ft.RotationZ))

	}

	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragTrack) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	boneCount := dec.Uint32()
	for i := 0; i < int(boneCount); i++ {
		ft := WldFragTrackBoneTransform{}
		if e.Flags&0x08 == 0x08 {
			ft.RotationDenominator = dec.Int16()
			ft.RotationX = dec.Int16()
			ft.RotationY = dec.Int16()
			ft.RotationZ = dec.Int16()
			ft.TranslationX = dec.Int16()
			ft.TranslationY = dec.Int16()
			ft.TranslationZ = dec.Int16()
			ft.TranslationDenominator = dec.Int16()
			e.BoneTransform = append(e.BoneTransform, ft)
			continue
		}
		ft.TranslationDenominator = int16(dec.Int8())
		ft.TranslationX = int16(dec.Int8())
		ft.TranslationY = int16(dec.Int8())
		ft.TranslationZ = int16(dec.Int8())
		ft.RotationDenominator = int16(dec.Int8())
		ft.RotationX = int16(dec.Int8())
		ft.RotationY = int16(dec.Int8())
		ft.RotationZ = int16(dec.Int8())
		e.BoneTransform = append(e.BoneTransform, ft)

	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// WldFragTrackRef is a bone in a skeleton. It is Track in libeq, Mob Skeleton Piece Track Reference in openzone, TRACKINSTANCE in wld, TrackDefFragment in lantern
type WldFragTrackRef struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	TrackRef int32  `yaml:"track_ref"`
	Flags    uint32 `yaml:"flags"`
	Sleep    uint32 `yaml:"sleep"` // if 0x01 is set, this is the number of milliseconds to sleep before starting the animation
}

func (e *WldFragTrackRef) FragCode() int {
	return 0x13
}

func (e *WldFragTrackRef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.TrackRef)
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

func (e *WldFragTrackRef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.TrackRef = dec.Int32()
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
type WldFragDMTrack struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragDMTrack) FragCode() int {
	return 0x2E
}

func (e *WldFragDMTrack) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDMTrack) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	return nil
}

// WldFragDMTrackRef is DmTrack in libeq, Mesh Animated Vertices Reference in openzone, empty in wld, MeshAnimatedVerticesReference in lantern
type WldFragDMTrackRef struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragDMTrackRef) FragCode() int {
	return 0x2F
}

func (e *WldFragDMTrackRef) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDMTrackRef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	return nil
}

// WldFragDMRGBTrack is a list of colors, one per vertex, for baked lighting. It is DmRGBTrackDef in libeq, Vertex Color in openzone, empty in wld, VertexColors in lantern
type WldFragDMRGBTrack struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragDMRGBTrack) FragCode() int {
	return 0x32
}

func (e *WldFragDMRGBTrack) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDMRGBTrack) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	return nil
}

// WldFragDMRGBTrackRef is DmRGBTrack in libeq, Vertex Color Reference in openzone, empty in wld, VertexColorsReference in lantern
type WldFragDMRGBTrackRef struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragDMRGBTrackRef) FragCode() int {
	return 0x33
}

func (e *WldFragDMRGBTrackRef) Write(w io.Writer) error {
	return nil
}

func (e *WldFragDMRGBTrackRef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	return nil
}
