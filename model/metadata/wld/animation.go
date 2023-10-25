package wld

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// SkeletonTrack is HierarchialSpriteDef in libeq, SkeletonTrackSet in openzone, HIERARCHIALSPRITE in wld, SkeletonHierarchy in lantern
type SkeletonTrack struct {
	NameRef            int32          `yaml:"name_ref"`
	Flags              uint32         `yaml:"flags"`
	AnimCount          uint32         `yaml:"anim_count"`
	CollisionVolumeRef uint32         `yaml:"collision_volume_ref"`
	CenterOffset       common.Vector3 `yaml:"center_offset"`
	Radius             float32        `yaml:"radius"`
}

func (e *SkeletonTrack) FragCode() int {
	return 0x10
}

func (e *SkeletonTrack) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.AnimCount)
	enc.Uint32(e.CollisionVolumeRef)
	if e.Flags&0x1 != 0 {
		enc.Float32(e.CenterOffset.X)
		enc.Float32(e.CenterOffset.Y)
		enc.Float32(e.CenterOffset.Z)
	}

	if e.Flags&0x2 != 0 {
		enc.Float32(e.Radius)
	}

	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeSkeletonTrack(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &SkeletonTrack{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.AnimCount = dec.Uint32()
	d.CollisionVolumeRef = dec.Uint32()
	if d.Flags&0x1 != 0 {
		d.CenterOffset.X = dec.Float32()
		d.CenterOffset.Y = dec.Float32()
		d.CenterOffset.Z = dec.Float32()
	}

	if d.Flags&0x2 != 0 {
		d.Radius = dec.Float32()
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// SkeletonTrackRef is HierarchialSprite in libeq, SkeletonTrackSetReference in openzone, HIERARCHIALSPRITE (ref) in wld, SkeletonHierarchyReference in lantern
type SkeletonTrackRef struct {
	NameRef          int16
	SkeletonTrackRef int16
	Flags            uint32
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

func decodeSkeletonTrackRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &SkeletonTrackRef{}
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
	NameRef            int32           `yaml:"name_ref"`
	Flags              uint32          `yaml:"flags"`
	SkeletonTransforms []BoneTransform `yaml:"skeleton_transforms"`
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
	enc.Uint32(uint32(len(e.SkeletonTransforms)))
	for _, ft := range e.SkeletonTransforms {
		enc.Int16(int16(ft.RotationDenominator))
		enc.Int16(int16(ft.RotationX))
		enc.Int16(int16(ft.RotationY))
		enc.Int16(int16(ft.RotationZ))
		enc.Int16(int16(ft.TranslationDenominator))
		enc.Int16(int16(ft.TranslationX))
		enc.Int16(int16(ft.TranslationY))
		enc.Int16(int16(ft.TranslationZ))
	}

	if enc.Error() != nil {
		return enc.Error()
	}

	return nil
}

func decodeTrack(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Track{}

	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	skeletonCount := dec.Uint32()
	for i := 0; i < int(skeletonCount); i++ {
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
		d.SkeletonTransforms = append(d.SkeletonTransforms, ft)

	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}

	return d, nil
}

// TrackRef is a bone in a skeleton. It is Track in libeq, Mob Skeleton Piece Track Reference in openzone, TRACKINSTANCE in wld, TrackDefFragment in lantern
type TrackRef struct {
	NameRef  int32
	TrackRef int32
	Flags    uint32
	Sleep    uint32 // if 0x01 is set, this is the number of milliseconds to sleep before starting the animation
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

func decodeTrackRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &TrackRef{}

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
