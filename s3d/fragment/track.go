package fragment

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/g3n/engine/math32"
)

// Track information
type Track struct {
	hashIndex uint32
	Frames    []*BoneTransform
}

// BoneTransform coordinate data
type BoneTransform struct {
	Translation math32.Vector3
	Rotation    math32.Quaternion
	Scale       float32
	ModelMatrix math32.Matrix4
}

func loadTrack(r io.ReadSeeker) (Fragment, error) {
	v := &Track{}
	err := parseTrack(r, v)
	if err != nil {
		return nil, fmt.Errorf("parse track: %w", err)
	}
	return v, nil
}

func parseTrack(r io.ReadSeeker, v *Track) error {
	if v == nil {
		return fmt.Errorf("track  is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &v.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flag: %w", err)
	}

	//if value &3 == 3 { // has data2 values
	//TODO: discern this
	//}

	if value != 8 { //object animations = 8
		return fmt.Errorf("unknown track type: expected 8, got %d", value)
	}

	var frameCount int
	err = binary.Read(r, binary.LittleEndian, &frameCount)
	if err != nil {
		return fmt.Errorf("read flag: %w", err)
	}
	for i := 0; i < frameCount; i++ {
		var rotDenominator, shiftDenominator int
		var rotX, rotY, rotZ int
		var shiftX, shiftY, shiftZ int

		err = binary.Read(r, binary.LittleEndian, &rotDenominator)
		if err != nil {
			return fmt.Errorf("read rot denominator: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, &rotX)
		if err != nil {
			return fmt.Errorf("read rot x: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, &rotY)
		if err != nil {
			return fmt.Errorf("read rot y: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, &rotZ)
		if err != nil {
			return fmt.Errorf("read rot z: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, &shiftX)
		if err != nil {
			return fmt.Errorf("read shift x: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, &shiftY)
		if err != nil {
			return fmt.Errorf("read shift y: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, &shiftZ)
		if err != nil {
			return fmt.Errorf("read shift z: %w", err)
		}

		err = binary.Read(r, binary.LittleEndian, &shiftDenominator)
		if err != nil {
			return fmt.Errorf("read shift denominator: %w", err)
		}

		ft := &BoneTransform{}

		if shiftDenominator != 0 {
			ft.Scale = float32(shiftDenominator) / 256
			ft.Translation.X = float32(shiftX / 256)
			ft.Translation.Y = float32(shiftY / 256)
			ft.Translation.Z = float32(shiftZ / 256)
		}
		ft.Rotation.X = float32(rotX)
		ft.Rotation.Y = float32(rotY)
		ft.Rotation.Z = float32(rotZ)
		ft.Rotation.W = float32(rotDenominator)
		ft.Rotation.Normalize()
		v.Frames = append(v.Frames, ft)
	}

	return nil
}

func (v *Track) FragmentType() string {
	return "Track "
}
