package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragTrackDef is TrackDef in libeq, Mob Skeleton Piece WldFragTrackDef in openzone, TRACKDEFINITION in wld, TrackDefFragment in lantern
type WldFragTrackDef struct {
	NameRef        int32                       `yaml:"name_ref"`
	Flags          uint32                      `yaml:"flags"`
	BoneTransforms []WldFragTrackBoneTransform `yaml:"skeleton_transforms"`
}

type WldFragTrackBoneTransform struct {
	RotateDenominator int16
	Rotation          [3]int16
	ShiftDenominator  int16
	Shift             [3]int16
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
			for _, shift := range ft.Shift {
				enc.Int16(shift)
			}
			for _, rotate := range ft.Rotation {
				enc.Int16(rotate)
			}
			enc.Int16(ft.RotateDenominator)
			continue
		}
		enc.Int8(int8(ft.RotateDenominator))
		for _, rotate := range ft.Rotation {
			enc.Int8(int8(rotate))
		}
		enc.Int8(int8(ft.ShiftDenominator))
		for _, shift := range ft.Shift {
			enc.Int8(int8(shift))
		}
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
			ft.Shift[0] = dec.Int16()
			ft.Shift[1] = dec.Int16()
			ft.Shift[2] = dec.Int16()
			ft.Rotation[0] = dec.Int16()
			ft.Rotation[1] = dec.Int16()
			ft.Rotation[2] = dec.Int16()
			ft.RotateDenominator = dec.Int16()
			e.BoneTransforms = append(e.BoneTransforms, ft)
			continue
		}
		ft.RotateDenominator = int16(dec.Int8())
		ft.Rotation[0] = int16(dec.Int8())
		ft.Rotation[1] = int16(dec.Int8())
		ft.Rotation[2] = int16(dec.Int8())
		ft.ShiftDenominator = int16(dec.Int8())
		ft.Shift[0] = int16(dec.Int8())
		ft.Shift[1] = int16(dec.Int8())
		ft.Shift[2] = int16(dec.Int8())
		e.BoneTransforms = append(e.BoneTransforms, ft)

	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
