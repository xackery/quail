package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/log"
)

type decoder struct {
	name  string
	parse func(r io.ReadSeeker, fragOffset int) error
}

// Decode loads a wld file
func (e *WLD) Decode(r io.ReadSeeker) error {

	var err error
	if e == nil {
		return fmt.Errorf("wld nil")
	}

	fragmentCount, err := e.readHeader(r)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	parsers := e.initParsers()

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(fragmentCount); fragOffset++ {
		var fragSize uint32
		var fragCode int32

		err = binary.Read(r, binary.LittleEndian, &fragSize)
		if err != nil {
			return fmt.Errorf("read fragment size %d/%d: %w", fragOffset, fragmentCount, err)
		}
		totalFragSize += fragSize
		//dump.Hex(fragSize, "%d(%s)fragSize=%d", i, name, fragSize)
		err = binary.Read(r, binary.LittleEndian, &fragCode)
		if err != nil {
			return fmt.Errorf("read fragment index %d/%d: %w", fragOffset, fragmentCount, err)
		}
		//dump.Hex(fragSize, "%dfragCode=%d", i, fragCode)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("frag position seek %d/%d: %w", fragOffset, fragmentCount, err)
		}

		buf := make([]byte, fragSize)
		_, err = r.Read(buf)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		parser, ok := parsers[fragCode]
		if !ok {
			log.Warnf("warning: unknown fragCode %d at offset %d", fragCode, fragOffset)
		} else {
			err = parser.parse(bytes.NewReader(buf), fragOffset)
			if err != nil {
				log.Warnf("warning: parse %s (%d, 0x%x): %s\n", parser.name, fragCode, fragCode, err)
				//return fmt.Errorf("parse %s: %w", parser.name, err)
			}
		}

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return fmt.Errorf("seek end of frag %d/%d: %w", fragOffset, fragmentCount, err)
		}
	}
	//dump.HexRange([]byte{byte(i), byte(i) + 1}, int(fragSize), "%dfrag=%s", i, frag.FragmentType())
	dump.HexRange([]byte{0, 1}, int(totalFragSize), "fragChunk=(%d bytes, %d entries)", int(totalFragSize), fragmentCount)

	for i, frag := range e.fragments {
		err = frag.build(e)
		if err != nil {
			return fmt.Errorf("build %d: %w", i, err)
		}
	}

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

func (e *WLD) readHeader(r io.ReadSeeker) (fragmentCount uint32, err error) {
	var header [4]byte
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		err = fmt.Errorf("read header: %w", err)
		return
	}

	validHeader := [4]byte{0x02, 0x3D, 0x50, 0x54}
	if header != validHeader {
		err = fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
		return
	}
	dump.Hex(header, "header=0x%x", header)

	err = binary.Read(r, binary.LittleEndian, &e.version)
	if err != nil {
		err = fmt.Errorf("read identifier: %w", err)
		return
	}

	e.isOldWorld = false
	switch e.version {
	case 0x00015500:
		e.isOldWorld = true
	case 0x1000C800:
		e.isOldWorld = false
	default:
		err = fmt.Errorf("unknown wld identifier %d", e.version)
		return
	}
	dump.Hex(e.version, "identifier=(isOld:%t)", e.isOldWorld)

	fragmentCount = uint32(0)
	err = binary.Read(r, binary.LittleEndian, &fragmentCount)
	if err != nil {
		err = fmt.Errorf("read fragmentCount: %w", err)
		return
	}
	dump.Hex(fragmentCount, "fragmentCount=%d", fragmentCount)

	unk1 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk1)
	if err != nil {
		err = fmt.Errorf("read unk1: %w", err)
		return
	}
	dump.Hex(unk1, "unk1=%d", unk1)

	unk2 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk2)
	if err != nil {
		err = fmt.Errorf("read unk2: %w", err)
		return
	}
	dump.Hex(unk2, "unk2=%d", unk2)

	var hashSize uint32
	err = binary.Read(r, binary.LittleEndian, &hashSize)
	if err != nil {
		err = fmt.Errorf("read hash size: %w", err)
		return
	}
	dump.Hex(hashSize, "hashSize=%d", hashSize)

	unk3 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk3)
	if err != nil {
		err = fmt.Errorf("read unk3: %w", err)
		return
	}
	dump.Hex(unk3, "unk3=%d", unk3)

	hashRaw := make([]byte, hashSize)

	err = binary.Read(r, binary.LittleEndian, &hashRaw)
	if err != nil {
		err = fmt.Errorf("read nameData: %w", err)
		return
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

	e.names = names

	dump.HexRange(hashRaw, int(hashSize), "nameData=(%d bytes, %d names)", hashSize, len(names))
	return
}
