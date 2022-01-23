package wld

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Save saves a wld file
func (e *WLD) Save(w io.Writer) error {
	if e == nil {
		return fmt.Errorf("wld nil")
	}
	var value uint32

	value = 0x54503D02
	err := binary.Write(w, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("write wld header: %w", err)
	}

	value = 0x00015500 //old world identifier
	err = binary.Write(w, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("write identifier: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, len(e.fragments))
	if err != nil {
		return fmt.Errorf("write fragment count: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, &e.BspRegionCount)
	if err != nil {
		return fmt.Errorf("write bsp region count: %w", err)
	}

	value = 0x680D4
	err = binary.Write(w, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("write after bsp region offset: %w", err)
	}

	var hashSize uint32
	err = binary.Write(w, binary.LittleEndian, &hashSize)
	if err != nil {
		return fmt.Errorf("write hash size: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("write after hash size offset: %w", err)
	}

	/*helper.WriteFixedString(w io.Writer, in string, size int)

	hashRaw, err := helper.ParseFixedString(r, hashSize)
	if err != nil {
		return fmt.Errorf("write hash: %w", err)
	}

	hashSplit := strings.Split(hashRaw, "\x00")
	e.Hash = make(map[int]string)
	offset := 0
	for _, h := range hashSplit {
		e.Hash[offset] = h
		offset += len(h) + 1
	}
	log.Debugf("fragments: %d, bsp regions: %d, bsp region offset: 0x%x, hashSize: %d", e.FragmentCount, e.BspRegionCount, value, hashSize)

	for i := 0; i < int(e.FragmentCount); i++ {
		var fragSize uint32
		var fragIndex int32

		err = binary.Write(w, binary.LittleEndian, &fragSize)
		if err != nil {
			return fmt.Errorf("write fragment size %d/%d: %w", i, e.FragmentCount, err)
		}
		err = binary.Write(w, binary.LittleEndian, &fragIndex)
		if err != nil {
			return fmt.Errorf("write fragment index %d/%d: %w", i, e.FragmentCount, err)
		}

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("frag position seek %d/%d: %w", i, e.FragmentCount, err)
		}

		buf := make([]byte, fragSize)
		_, err = r.Write(buf)
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}

		log.Debugf("%d fragIndex: %d 0x%x, len %d\n%s", i, fragIndex, fragIndex, len(buf), hex.Dump(buf))
		frag, err := e.ParseFragment(fragIndex, bytes.NewWriteer(buf))
		if err != nil {
			return fmt.Errorf("fragment load: %w", err)
		}
		log.Debugf("%d fragIndex: %d 0x%x determined to be %s", i, fragIndex, fragIndex, frag.FragmentType())

		e.Fragments = append(e.Fragments, frag)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return fmt.Errorf("seek end of frag %d/%d: %w", i, e.FragmentCount, err)
		}
	}
	*/
	return nil
}
