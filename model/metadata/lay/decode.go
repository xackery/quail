package lay

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
)

func (e *LAY) Decode(r io.ReadSeeker) error {
	var err error
	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'E', 'Q', 'G', 'L'} {
		return fmt.Errorf("header does not match EQGL")
	}

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("read version: %w", err)
	}
	dump.Hex(version, "version=%d", version)
	versionOffset := 0
	switch version {
	case 2:
		versionOffset = 40 //32
	case 3:
		versionOffset = 18 //14
	case 4:
		versionOffset = 20
	default:

	}

	nameLength := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &nameLength)
	if err != nil {
		return fmt.Errorf("read nameLength: %w", err)
	}
	dump.Hex(nameLength, "nameLength=%d", nameLength)

	materialCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &materialCount)
	if err != nil {
		return fmt.Errorf("read materialCount: %w", err)
	}
	dump.Hex(materialCount, "materialCount=%d", materialCount)

	nameData := make([]byte, nameLength)

	err = binary.Read(r, binary.LittleEndian, &nameData)
	if err != nil {
		return fmt.Errorf("read nameData: %w", err)
	}

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

	dump.HexRange(nameData, int(nameLength), "nameData=(%d bytes, %d entries)", nameLength, len(names))
	for i := 0; i < int(materialCount); i++ {
		materialID := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read materialID: %w", err)
		}

		name, ok := names[materialID]
		if !ok {
			return fmt.Errorf("%d names materialID 0x%x not found", i, materialID)
		}

		dump.Hex(materialID, "%dmaterialID=%d(%s)", i, materialID, name)

		entry0Offset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &entry0Offset)
		if err != nil {
			return fmt.Errorf("read entry0ID: %w", err)
		}
		entry0Name, ok := names[entry0Offset]
		if !ok {
			return fmt.Errorf("%d names entry0Offset 0x%x not found", i, entry0Offset)
		}
		dump.Hex(entry0Offset, "%dentry0ID=0x%x(%s)", i, entry0Offset, entry0Name)

		entry1Offset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &entry1Offset)
		if err != nil {
			return fmt.Errorf("read entry1ID: %w", err)
		}

		entry1Name := ""

		if entry1Offset != 0xffffffff {
			entry1Name, ok = names[entry1Offset]
			if !ok {
				return fmt.Errorf("%dnames entry1Offset 0x%x not found", i, entry1Offset)
			}
		}
		dump.Hex(entry1Offset, "%dentry1ID=0x%x(%s)", i, entry1Offset, entry1Name)

		err = e.MaterialAdd(name, entry0Name, entry1Name)
		if err != nil {
			return fmt.Errorf("materialAdd: %w", err)
		}

		_, err = r.Seek(int64(versionOffset), io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("%dseek version %d: %w", i, version, err)
		}
	}

	return nil
}
