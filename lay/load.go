package lay

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
)

func (e *LAY) Load(r io.ReadSeeker) error {
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
		versionOffset = 32
	case 3:
		versionOffset = 14
	case 4:
		versionOffset = 16
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

	fmt.Println(hex.Dump(nameData))
	dump.HexRange(nameData, int(nameLength), "nameData=(%d bytes, %d entries)", nameLength, len(names))
	for i := 0; i < int(materialCount); i++ {
		materialID := int32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read materialID: %w", err)
		}
		dump.Hex(materialID, "%dmaterialID=%d", i, materialID)

		diffuseOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &diffuseOffset)
		if err != nil {
			return fmt.Errorf("read diffuseID: %w", err)
		}
		diffuseName, ok := names[diffuseOffset]
		if !ok {
			return fmt.Errorf("%d names diffuseOffset 0x%x not found", i, diffuseOffset)
		}
		dump.Hex(diffuseOffset, "%ddiffuseID=0x%x(%s)", i, diffuseOffset, diffuseName)

		normalOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &normalOffset)
		if err != nil {
			return fmt.Errorf("read normalID: %w", err)
		}
		normalName, ok := names[normalOffset]
		if !ok {
			return fmt.Errorf("%d names normal offset 0x%x not found", i, normalOffset)
		}
		dump.Hex(normalOffset, "%dnormalID=0x%x(%s)", i, normalOffset, normalName)

		err = e.MaterialAdd(diffuseName, normalName)
		if err != nil {
			return fmt.Errorf("materialADD: %w", err)
		}

		_, err = r.Seek(int64(versionOffset), io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("%dseek version %d: %w", i, version, err)
		}
	}

	return nil
}
