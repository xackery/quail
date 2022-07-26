package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// Track information
type Track struct {
	name   string
	Frames []*BoneTransform
}

// BoneTransform coordinate data
type BoneTransform struct {
	Translation [3]float32
	Rotation    [4]float32
	Scale       float32
	ModelMatrix [16]float32
}

func LoadTrack(r io.ReadSeeker) (common.WldFragmenter, error) {
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
	var err error
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
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

	var frameCount int16
	err = binary.Read(r, binary.LittleEndian, &frameCount)
	if err != nil {
		return fmt.Errorf("read frameCount: %w", err)
	}
	for i := 0; i < int(frameCount); i++ {
		var rotDenominator, shiftDenominator int16
		var rotX, rotY, rotZ int16
		var shiftX, shiftY, shiftZ int16

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
			ft.Translation[0] = float32(shiftX / 256)
			ft.Translation[1] = float32(shiftY / 256)
			ft.Translation[2] = float32(shiftZ / 256)
		}
		ft.Rotation[0] = float32(rotX)
		ft.Rotation[1] = float32(rotY)
		ft.Rotation[2] = float32(rotZ)
		ft.Rotation[3] = float32(rotDenominator)
		ft.Rotation = helper.Normalize(ft.Rotation)
		v.Frames = append(v.Frames, ft)
	}

	return nil
}

func (v *Track) FragmentType() string {
	return "Track "
}

func (e *Track) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
