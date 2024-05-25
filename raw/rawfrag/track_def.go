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
