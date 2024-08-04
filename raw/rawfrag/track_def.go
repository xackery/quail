package rawfrag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// WldFragTrackDef is TrackDef in libeq, Mob Skeleton Piece WldFragTrackDef in openzone, TRACKDEFINITION in wld, TrackDefFragment in lantern
type WldFragTrackDef struct {
	NameRef         int32                       `yaml:"name_ref"`
	Flags           uint32                      `yaml:"flags"`
	FrameTransforms []WldFragTrackBoneTransform `yaml:"skeleton_transforms"`
}

type WldFragTrackBoneTransform struct {
	RotateDenominator int16
	Rotation          [4]int16
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
	enc.Uint32(uint32(len(e.FrameTransforms)))
	for _, ft := range e.FrameTransforms {
		if e.Flags&0x08 == 0x08 {
			enc.Int8(int8(ft.RotateDenominator))
			enc.Int8(int8(ft.Rotation[0]))
			enc.Int8(int8(ft.Rotation[1]))
			enc.Int8(int8(ft.Rotation[2]))
			enc.Int8(int8(ft.ShiftDenominator))
			enc.Int8(int8(ft.Shift[0]))
			enc.Int8(int8(ft.Shift[1]))
			enc.Int8(int8(ft.Shift[2]))
			continue
		}
		enc.Int16(ft.ShiftDenominator)
		enc.Int16(ft.Shift[0])
		enc.Int16(ft.Shift[1])
		enc.Int16(ft.Shift[2])
		enc.Int16(ft.Rotation[0])
		enc.Int16(ft.Rotation[1])
		enc.Int16(ft.Rotation[2])
		enc.Int16(ft.RotateDenominator)
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
	frameCount := dec.Uint32()
	for i := 0; i < int(frameCount); i++ {
		ft := WldFragTrackBoneTransform{}
		if e.Flags&0x08 == 0x08 {
			ft.RotateDenominator = int16(dec.Int8())
			ft.Rotation[0] = int16(dec.Int8())
			ft.Rotation[1] = int16(dec.Int8())
			ft.Rotation[2] = int16(dec.Int8())
			ft.ShiftDenominator = int16(dec.Int8())
			ft.Shift[0] = int16(dec.Int8())
			ft.Shift[1] = int16(dec.Int8())
			ft.Shift[2] = int16(dec.Int8())
			e.FrameTransforms = append(e.FrameTransforms, ft)
			continue
		}
		ft.ShiftDenominator = dec.Int16()
		ft.Shift[0] = dec.Int16()
		ft.Shift[1] = dec.Int16()
		ft.Shift[2] = dec.Int16()
		ft.Rotation[0] = dec.Int16()
		ft.Rotation[1] = dec.Int16()
		ft.Rotation[2] = dec.Int16()
		ft.RotateDenominator = dec.Int16()
		e.FrameTransforms = append(e.FrameTransforms, ft)
	}

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
