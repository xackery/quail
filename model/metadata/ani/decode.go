package ani

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/model/geo"
)

func (e *ANI) Decode(r io.ReadSeeker) error {
	var err error
	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'E', 'Q', 'G', 'A'} {
		return fmt.Errorf("header does not match EQGM")
	}

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("read version: %w", err)
	}
	dump.Hex(version, "version=%d", version)

	nameLength := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &nameLength)
	if err != nil {
		return fmt.Errorf("read nameLength: %w", err)
	}
	dump.Hex(nameLength, "nameLength=%d", nameLength)

	boneCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &boneCount)
	if err != nil {
		return fmt.Errorf("read boneCount: %w", err)
	}
	dump.Hex(boneCount, "boneCount=%d", boneCount)

	e.isStrict = false
	if version > 1 {
		isStrict := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &isStrict)
		if err != nil {
			return fmt.Errorf("read isStrict: %w", err)
		}
		dump.Hex(isStrict, "isStrict=%d", isStrict)

		if isStrict > 0 {
			e.isStrict = true
		}
	}

	nameData := make([]byte, nameLength)
	err = binary.Read(r, binary.LittleEndian, &nameData)
	if err != nil {
		return fmt.Errorf("read nameData: %w", err)
	}
	dump.Hex(nameData, "nameData=(%d bytes)", len(nameData))

	names := make(map[uint32]string)
	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	for i := 0; i < int(boneCount); i++ {
		bone := &geo.BoneAnimation{}

		err = binary.Read(r, binary.LittleEndian, &bone.FrameCount)
		if err != nil {
			return fmt.Errorf("read bone %d frameCount: %w", i, err)
		}
		dump.Hex(bone.FrameCount, "%dframeCount %d", i, bone.FrameCount)

		nameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &nameOffset)
		if err != nil {
			return fmt.Errorf("read bone %d nameOffset: %w", i, err)
		}
		bone.Name = names[nameOffset]
		dump.Hex(nameOffset, "%dnameOffset 0x%x (%s)", i, nameOffset, bone.Name)

		for j := 0; j < int(bone.FrameCount); j++ {
			frame := &geo.BoneAnimationFrame{}

			err = binary.Read(r, binary.LittleEndian, &frame.Milliseconds)
			if err != nil {
				return fmt.Errorf("read bone%d frame%d frame.milliseconds: %w", i, j, err)
			}

			translation := &geo.Vector3{}
			err = binary.Read(r, binary.LittleEndian, translation)
			if err != nil {
				return fmt.Errorf("read bone%d frame%d frame.translation: %w", i, j, err)
			}
			frame.Translation = translation

			rotation := &geo.Quad4{}
			err = binary.Read(r, binary.LittleEndian, rotation)
			if err != nil {
				return fmt.Errorf("read bone%d frame%d frame.rotation: %w", i, j, err)
			}
			frame.Rotation = rotation

			scale := &geo.Vector3{}
			err = binary.Read(r, binary.LittleEndian, scale)
			if err != nil {
				return fmt.Errorf("read bone%d frame%d frame.scale: %w", i, j, err)
			}
			frame.Scale = scale
			bone.Frames = append(bone.Frames, frame)
		}

		dump.HexRange([]byte{0x01, 0x02}, int(bone.FrameCount*44), "%dbone frames(%d bytes)", i, bone.FrameCount*44)
		e.bones = append(e.bones, bone)
	}

	return nil
}
