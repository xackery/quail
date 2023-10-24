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
	NameRef            int32
	Flags              uint32
	SkeletonTransforms []BoneTransform
}

type BoneTransform struct {
	Translation common.Vector3
	Rotation    common.Quad4
	Scale       float32
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
		enc.Int16(int16(ft.Rotation.W))
		enc.Int16(int16(ft.Rotation.X))
		enc.Int16(int16(ft.Rotation.Y))
		enc.Int16(int16(ft.Rotation.Z))
		enc.Int16(int16(ft.Translation.X))
		enc.Int16(int16(ft.Translation.Y))
		enc.Int16(int16(ft.Translation.Z))
		enc.Int16(int16(ft.Scale * 256))
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
		rotDenom := dec.Int16()
		rotX := dec.Int16()
		rotY := dec.Int16()
		rotZ := dec.Int16()
		shiftX := dec.Int16()
		shiftY := dec.Int16()
		shiftZ := dec.Int16()
		shiftDenom := dec.Int16()
		ft := BoneTransform{}
		if shiftDenom != 0 {
			ft.Scale = float32(shiftDenom) / 256
			ft.Translation.X = float32(shiftX) / 256
			ft.Translation.Y = float32(shiftY) / 256
			ft.Translation.Z = float32(shiftZ) / 256
		}
		ft.Rotation.X = float32(rotX)
		ft.Rotation.Y = float32(rotY)
		ft.Rotation.Z = float32(rotZ)
		ft.Rotation.W = float32(rotDenom)
		ft.Rotation = common.Normalize(ft.Rotation)
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
