package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/model/mesh/wld/fragment"
)

// Decode loads a wld file
func (e *WLD) Decode(r io.ReadSeeker) error {
	var err error
	if e == nil {
		return fmt.Errorf("wld nil")
	}
	var value uint32

	var header [4]byte
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	validHeader := [4]byte{0x02, 0x3D, 0x50, 0x54}
	if header != validHeader {
		return fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	dump.Hex(header, "header=0x%x", header)

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("read identifier: %w", err)
	}

	isOldWorld := false
	switch version {
	case 0x00015500:
		isOldWorld = true
	case 0x1000C800:
		isOldWorld = false
	default:
		return fmt.Errorf("unknown wld identifier %d", value)
	}
	dump.Hex(value, "identifier=(isOld:%t)", isOldWorld)

	fragmentCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &fragmentCount)
	if err != nil {
		return fmt.Errorf("read fragmentCount: %w", err)
	}
	dump.Hex(fragmentCount, "fragmentCount=%d", fragmentCount)

	unk1 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk1)
	if err != nil {
		return fmt.Errorf("read unk1: %w", err)
	}
	dump.Hex(unk1, "unk1=%d", unk1)

	unk2 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk2)
	if err != nil {
		return fmt.Errorf("read unk2: %w", err)
	}
	dump.Hex(unk2, "unk2=%d", unk2)

	var hashSize uint32
	err = binary.Read(r, binary.LittleEndian, &hashSize)
	if err != nil {
		return fmt.Errorf("read hash size: %w", err)
	}
	dump.Hex(hashSize, "hashSize=%d", hashSize)

	unk3 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk3)
	if err != nil {
		return fmt.Errorf("read unk3: %w", err)
	}
	dump.Hex(unk3, "unk3=%d", unk3)

	hashRaw := make([]byte, hashSize)

	err = binary.Read(r, binary.LittleEndian, &hashRaw)
	if err != nil {
		return fmt.Errorf("read nameData: %w", err)
	}

	nameData := decodeStringHash(hashRaw)

	names := make(map[int32]string)

	chunk := []rune{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[int32(lastOffset)] = string(chunk)
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	fragment.SetNames(names)

	e.NameCache = names

	dump.HexRange(hashRaw, int(hashSize), "nameData=(%d bytes, %d names)", hashSize, len(names))

	parsers := []struct {
		invoke func(frag *fragmentInfo) error
		name   string
	}{
		{invoke: e.parseMesh, name: "mesh"},
		{invoke: e.parseMaterial, name: "material"},
	}

	totalFragSize := uint32(0)
	for i := 0; i < int(fragmentCount); i++ {
		var fragSize uint32
		var fragIndex int32

		name := names[int32(i)]
		err = binary.Read(r, binary.LittleEndian, &fragSize)
		if err != nil {
			return fmt.Errorf("read fragment size %d/%d: %w", i, fragmentCount, err)
		}
		totalFragSize += fragSize
		//dump.Hex(fragSize, "%d(%s)fragSize=%d", i, name, fragSize)
		err = binary.Read(r, binary.LittleEndian, &fragIndex)
		if err != nil {
			return fmt.Errorf("read fragment index %d/%d: %w", i, fragmentCount, err)
		}
		//dump.Hex(fragSize, "%dfragIndex=%d", i, fragIndex)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("frag position seek %d/%d: %w", i, fragmentCount, err)
		}
		if fragIndex == 0x03 {
			fmt.Println("fragPos for 0x03:", fragPosition)
		}

		buf := make([]byte, fragSize)
		_, err = r.Read(buf)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		frag, err := fragment.New(fragIndex, bytes.NewReader(buf))
		if err != nil {
			//TODO: fix error
			//fmt.Printf("warning: fragment decode 0x%x (%d): %s\n", fragIndex, fragIndex, err)
			//return fmt.Errorf("fragment load: %w", err)
		} else {
			for _, parser := range parsers {
				err = parser.invoke(&fragmentInfo{name: name, data: frag})
				if err != nil {
					//fmt.Printf("warning: parse %s: %s\n", parser.name, err)
					//return fmt.Errorf("parse %s: %w", parser.name, err)
				}
			}
		}

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return fmt.Errorf("seek end of frag %d/%d: %w", i, fragmentCount, err)
		}
	}
	//dump.HexRange([]byte{byte(i), byte(i) + 1}, int(fragSize), "%dfrag=%s", i, frag.FragmentType())
	dump.HexRange([]byte{0, 1}, int(totalFragSize), "fragChunk=(%d bytes, %d entries)", int(totalFragSize), fragmentCount)

	// Now convert fragments to data

	return nil
}

func decodeStringHash(hash []byte) string {
	hashKey := []byte{0x95, 0x3A, 0xC5, 0x2A, 0x95, 0x7A, 0x95, 0x6A}
	out := ""
	for i := 0; i < len(hash); i++ {
		out += string(hash[i] ^ hashKey[i%8])
	}
	return out
}
