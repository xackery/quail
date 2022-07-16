package ani

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/dump"
)

func (e *ANI) Load(r io.ReadSeeker) error {
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

	animationCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &animationCount)
	if err != nil {
		return fmt.Errorf("read animationCount: %w", err)
	}
	dump.Hex(animationCount, "animationCount=%d", animationCount)

	e.isStrict = false
	if version == 1 {
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
		}
		chunk = append(chunk, b)
	}

	for i := 0; i < int(animationCount); i++ {
		bone := &Bone{}

		err = binary.Read(r, binary.LittleEndian, &bone.frameCount)
		if err != nil {
			return fmt.Errorf("read bone %d frameCount: %w", i, err)
		}

		nameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &nameOffset)
		if err != nil {
			return fmt.Errorf("read bone %d nameOffset: %w", i, err)
		}

		bone.name = names[nameOffset]

		err = binary.Read(r, binary.LittleEndian, &bone.delay)
		if err != nil {
			return fmt.Errorf("read %d bone.delay: %w", i, err)
		}
		translation := math32.NewVec3()
		err = binary.Read(r, binary.LittleEndian, translation)
		if err != nil {
			return fmt.Errorf("read %d bone.translation: %w", i, err)
		}
		bone.translation = translation

		rotation := math32.NewVec4()
		err = binary.Read(r, binary.LittleEndian, rotation)
		if err != nil {
			return fmt.Errorf("read %d bone.rotation: %w", i, err)
		}
		bone.rotation = rotation

		scale := math32.NewVec3()
		err = binary.Read(r, binary.LittleEndian, scale)
		if err != nil {
			return fmt.Errorf("read %d bone.scale: %w", i, err)
		}
		bone.scale = scale

		dump.Hex(chunk, "%dbone(%d %+v)", i, bone.delay, bone.translation)
		e.bones = append(e.bones, bone)
	}

	return nil
}
