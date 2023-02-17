package ani

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
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
		bone := &Bone{}

		err = binary.Read(r, binary.LittleEndian, &bone.frameCount)
		if err != nil {
			return fmt.Errorf("read bone %d frameCount: %w", i, err)
		}
		dump.Hex(bone.frameCount, "%dframeCount %d", i, bone.frameCount)

		nameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &nameOffset)
		if err != nil {
			return fmt.Errorf("read bone %d nameOffset: %w", i, err)
		}
		bone.name = names[nameOffset]
		dump.Hex(nameOffset, "%dnameOffset 0x%x (%s)", i, nameOffset, bone.name)

		for j := 0; j < int(bone.frameCount); j++ {
			frame := &Frame{}

			err = binary.Read(r, binary.LittleEndian, &frame.milliseconds)
			if err != nil {
				return fmt.Errorf("read bone%d frame%d frame.milliseconds: %w", i, j, err)
			}

			translation := [3]float32{}
			err = binary.Read(r, binary.LittleEndian, &translation)
			if err != nil {
				return fmt.Errorf("read bone%d frame%d frame.translation: %w", i, j, err)
			}
			frame.translation = translation

			rotation := [4]float32{}
			err = binary.Read(r, binary.LittleEndian, &rotation)
			if err != nil {
				return fmt.Errorf("read bone%d frame%d frame.rotation: %w", i, j, err)
			}
			frame.rotation = rotation

			scale := [3]float32{}
			err = binary.Read(r, binary.LittleEndian, &scale)
			if err != nil {
				return fmt.Errorf("read bone%d frame%d frame.scale: %w", i, j, err)
			}
			frame.scale = scale
			bone.frames = append(bone.frames, frame)
		}

		dump.HexRange([]byte{0x01, 0x02}, int(bone.frameCount*44), "%dbone frames(%d bytes)", i, bone.frameCount*44)
		e.bones = append(e.bones, bone)
	}

	return nil
}
