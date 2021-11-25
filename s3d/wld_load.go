package s3d

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/s3d/fragment"
)

// Load loads a wld file
func (e *Wld) Load(r io.ReadSeeker) error {
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

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read identifier: %w", err)
	}
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
		return fmt.Errorf("read fragment count: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &e.BspRegionCount)
	if err != nil {
		return fmt.Errorf("read bsp region count: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read after bsp region offset: %w", err)
	}
	//if value != 0x000680D4 {
	//	return fmt.Errorf("after bsp region offset wanted 0x%x, got 0x%x", 0x000680D4, value)
	//}//

	var hashSize uint32
	err = binary.Read(r, binary.LittleEndian, &hashSize)
	if err != nil {
		return fmt.Errorf("read hash size: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read after hash size offset: %w", err)
	}

	hashRaw, err := helper.ParseFixedString(r, hashSize)
	if err != nil {
		return fmt.Errorf("read hash: %w", err)
	}

	hashSplit := strings.Split(hashRaw, "\x00")
	e.Hash = make(map[int]string)
	offset := 0
	for _, h := range hashSplit {
		e.Hash[offset] = h
		offset += len(h) + 1
	}

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
		frag, err := fragment.Load(fragIndex, r)
		if err != nil {
			return fmt.Errorf("fragment load: %w", err)
		}
		e.Fragments = append(e.Fragments, frag)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return fmt.Errorf("seek end of frag %d/%d: %w", i, e.FragmentCount, err)
		}
	}
	return nil
}
