package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/dump"
)

// Load loads a wld file
func (e *WLD) Load(r io.ReadSeeker) error {
	if e == nil {
		return fmt.Errorf("wld nil")
	}
	var value uint32

	err := binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read wld header: %w", err)
	}
	if value != 0x54503D02 {
		return fmt.Errorf("unknown wld header: wanted 0x%x, got 0x%x", 0x54503D02, value)
	}
	dump.Hex(value, "header=%d", value)

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read identifier: %w", err)
	}
	dump.Hex(value, "identifier=0x%x", value)
	switch value {
	case 0x00015500:
		e.IsOldWorld = true
	case 0x1000C800:
		e.IsOldWorld = false
	default:
		return fmt.Errorf("unknown wld identifier %d", value)
	}

	err = binary.Read(r, binary.LittleEndian, &e.FragmentCount)
	if err != nil {
		return fmt.Errorf("read fragmentCount: %w", err)
	}
	dump.Hex(e.FragmentCount, "fragmentCount=%d", e.FragmentCount)

	err = binary.Read(r, binary.LittleEndian, &e.BspRegionCount)
	if err != nil {
		return fmt.Errorf("read BspRegionCount: %w", err)
	}
	dump.Hex(e.BspRegionCount, "BspRegionCount=%d", e.BspRegionCount)

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read after bsp region offset: %w", err)
	}
	if value != 0x680D4 {
		return fmt.Errorf("after bsp region offset wanted 0x%x, got 0x%x", 0x680D4, value)
	}

	dump.Hex(value, "BspRegionAfterOffset=0x%x", value)

	var hashSize uint32
	err = binary.Read(r, binary.LittleEndian, &hashSize)
	if err != nil {
		return fmt.Errorf("read hash size: %w", err)
	}
	dump.Hex(hashSize, "hashSize=%d", hashSize)

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read after hash size offset: %w", err)
	}
	dump.Hex(value, "afterHashOffset=0x%x", value)

	hashRaw := make([]byte, hashSize)

	err = binary.Read(r, binary.LittleEndian, hashRaw)
	if err != nil {
		return fmt.Errorf("read hash raw: %w", err)
	}
	dump.Hex(hashRaw, "hashRaw=(%d bytes)", len(hashRaw))

	hashString := decodeStringHash(hashRaw)

	hashSplit := strings.Split(hashString, "\x00")
	e.Hash = make(map[int]string)
	offset := 0
	for _, h := range hashSplit {
		e.Hash[offset] = h
		offset += len(h) + 1
		//log.Debugf("adding hash at 0x%x", offset)
	}
	//log.Debugf("fragments: %d, bsp regions: %d, bsp region offset: 0x%x, hashSize: %d", e.FragmentCount, e.BspRegionCount, value, hashSize)

	for i := 0; i < int(e.FragmentCount); i++ {
		var fragSize uint32
		var fragIndex int32

		err = binary.Read(r, binary.LittleEndian, &fragSize)
		if err != nil {
			return fmt.Errorf("read fragment size %d/%d: %w", i, e.FragmentCount, err)
		}
		err = binary.Read(r, binary.LittleEndian, &fragIndex)
		if err != nil {
			return fmt.Errorf("read fragment index %d/%d: %w", i, e.FragmentCount, err)
		}

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("frag position seek %d/%d: %w", i, e.FragmentCount, err)
		}

		buf := make([]byte, fragSize)
		_, err = r.Read(buf)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		//log.Debugf("%d fragIndex: %d 0x%x, len %d\n%s", i, fragIndex, fragIndex, len(buf), hex.Dump(buf))
		frag, err := e.ParseFragment(fragIndex, bytes.NewReader(buf))
		if err != nil {
			return fmt.Errorf("fragment load: %w", err)
		}
		//log.Debugf("%d fragIndex: %d 0x%x determined to be %s", i, fragIndex, fragIndex, frag.FragmentType())

		e.Fragments = append(e.Fragments, frag)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return fmt.Errorf("seek end of frag %d/%d: %w", i, e.FragmentCount, err)
		}
	}
	return nil
}

func decodeStringHash(hash []byte) string {
	hashKey := []byte{0x95, 0x3A, 0xC5, 0x2A, 0x95, 0x7A, 0x95, 0x6A}
	out := ""
	for i := 0; i < len(hash); i++ {
		out = string(hash[i] ^ hashKey[i%8])
	}
	return out
}
