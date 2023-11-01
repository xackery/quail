package raw

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
)

// SkeletonTrack is HierarchialSpriteDef in libeq, SkeletonTrackSet in openzone, HIERARCHIALSPRITE in wld, SkeletonHierarchy in lantern
type SkeletonTrack struct {
	FragName           string          `yaml:"frag_name"`
	NameRef            int32           `yaml:"name_ref"`
	Flags              uint32          `yaml:"flags"`
	CollisionVolumeRef uint32          `yaml:"collision_volume_ref"`
	CenterOffset       [3]uint32       `yaml:"center_offset"`
	BoundingRadius     float32         `yaml:"bounding_radius"`
	Bones              []SkeletonEntry `yaml:"bones"`
	Skins              []uint32        `yaml:"skins"`
	SkinLinks          []uint32        `yaml:"skin_links"`
}

type SkeletonEntry struct {
	NameRef         int32    `yaml:"name_ref"`
	Flags           uint32   `yaml:"flags"`
	TrackRef        uint32   `yaml:"track_ref"`
	MeshOrSpriteRef uint32   `yaml:"mesh_or_sprite_ref"`
	SubBones        []uint32 `yaml:"sub_bones"`
}

func (e *SkeletonTrack) FragCode() int {
	return 0x10
}

func (e *SkeletonTrack) Encode(w io.Writer) error {
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

	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func readSkeletonTrack(r io.ReadSeeker) (FragmentReader, error) {
	d := &SkeletonTrack{}
	d.FragName = FragName(d.FragCode())

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	boneCount := dec.Uint32()
	d.CollisionVolumeRef = dec.Uint32()
	if d.Flags&0x1 != 0 {
		d.CenterOffset = [3]uint32{dec.Uint32(), dec.Uint32(), dec.Uint32()}
	}

	if d.Flags&0x2 != 0 {
		d.BoundingRadius = dec.Float32()
	}

	for i := 0; i < int(boneCount); i++ {
		bone := SkeletonEntry{}
		bone.NameRef = dec.Int32()
		bone.Flags = dec.Uint32()
		bone.TrackRef = dec.Uint32()
		bone.MeshOrSpriteRef = dec.Uint32()
		subCount := dec.Uint32()
		for j := 0; j < int(subCount); j++ {
			bone.SubBones = append(bone.SubBones, dec.Uint32())
		}
		d.Bones = append(d.Bones, bone)
	}

	if d.Flags&0x200 != 0 {
		skinCount := dec.Uint32()
		for i := 0; i < int(skinCount); i++ {
			d.Skins = append(d.Skins, dec.Uint32())
		}
		for i := 0; i < int(skinCount); i++ {
			d.SkinLinks = append(d.SkinLinks, dec.Uint32())
		}
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// SkeletonTrackRef is HierarchialSprite in libeq, SkeletonTrackSetReference in openzone, HIERARCHIALSPRITE (ref) in wld, SkeletonHierarchyReference in lantern
type SkeletonTrackRef struct {
	FragName         string `yaml:"frag_name"`
	NameRef          int16  `yaml:"name_ref"`
	SkeletonTrackRef int16  `yaml:"skeleton_track_ref"`
	Flags            uint32 `yaml:"flags"`
}

func (e *SkeletonTrackRef) FragCode() int {
	return 0x11
}

func (e *SkeletonTrackRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int16(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Int16(e.SkeletonTrackRef)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func readSkeletonTrackRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &SkeletonTrackRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int16()
	d.Flags = dec.Uint32()
	d.SkeletonTrackRef = dec.Int16()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// Track is TrackDef in libeq, Mob Skeleton Piece Track in openzone, TRACKDEFINITION in wld, TrackDefFragment in lantern
type Track struct {
	FragName      string          `yaml:"frag_name"`
	NameRef       int32           `yaml:"name_ref"`
	Flags         uint32          `yaml:"flags"`
	BoneTransform []BoneTransform `yaml:"skeleton_transforms"`
}

type BoneTransform struct {
	TranslationDenominator int16 `yaml:"translation_denominator"`
	TranslationX           int16 `yaml:"translation_x"`
	TranslationY           int16 `yaml:"translation_y"`
	TranslationZ           int16 `yaml:"translation_z"`
	RotationDenominator    int16 `yaml:"rotation_denominator"`
	RotationX              int16 `yaml:"rotation_x"`
	RotationY              int16 `yaml:"rotation_y"`
	RotationZ              int16 `yaml:"rotation_z"`
}

func (e *Track) FragCode() int {
	return 0x12
}

func (e *Track) Encode(w io.Writer) error {
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

	if enc.Error() != nil {
		return enc.Error()
	}

	return nil
}

func readTrack(r io.ReadSeeker) (FragmentReader, error) {
	d := &Track{}
	d.FragName = FragName(d.FragCode())

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	boneCount := dec.Uint32()
	for i := 0; i < int(boneCount); i++ {
		ft := BoneTransform{}
		if d.Flags&0x08 == 0x08 {
			ft.RotationDenominator = dec.Int16()
			ft.RotationX = dec.Int16()
			ft.RotationY = dec.Int16()
			ft.RotationZ = dec.Int16()
			ft.TranslationX = dec.Int16()
			ft.TranslationY = dec.Int16()
			ft.TranslationZ = dec.Int16()
			ft.TranslationDenominator = dec.Int16()
			d.BoneTransform = append(d.BoneTransform, ft)
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
		d.BoneTransform = append(d.BoneTransform, ft)

	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}

// TrackRef is a bone in a skeleton. It is Track in libeq, Mob Skeleton Piece Track Reference in openzone, TRACKINSTANCE in wld, TrackDefFragment in lantern
type TrackRef struct {
	FragName string `yaml:"frag_name"`
	NameRef  int32  `yaml:"name_ref"`
	TrackRef int32  `yaml:"track_ref"`
	Flags    uint32 `yaml:"flags"`
	Sleep    uint32 `yaml:"sleep"` // if 0x01 is set, this is the number of milliseconds to sleep before starting the animation
}

func (e *TrackRef) FragCode() int {
	return 0x13
}

func (e *TrackRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.TrackRef)
	enc.Uint32(e.Flags)
	if e.Flags&0x01 == 0x01 {
		enc.Uint32(e.Sleep)
	}

	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func readTrackRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &TrackRef{}
	d.FragName = FragName(d.FragCode())

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.TrackRef = dec.Int32()
	d.Flags = dec.Uint32()
	if d.Flags&0x01 == 0x01 {
		d.Sleep = dec.Uint32()
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}
