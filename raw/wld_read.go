package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/tag"
)

type Wld struct {
	MetaFileName   string                     `yaml:"file_name"`
	Version        uint32                     `yaml:"version"`
	IsOldWorld     bool                       `yaml:"is_old_world"`
	Fragments      []model.FragmentReadWriter `yaml:"fragments,omitempty"`
	BspRegionCount uint32                     `yaml:"bsp_region_count"`
	Unk2           uint32                     `yaml:"unk2"`
	Unk3           uint32                     `yaml:"unk3"`
}

func (wld *Wld) Identity() string {
	return "wld"
}

// Read reads a wld file that was prepped by Load
func (wld *Wld) Read(r io.ReadSeeker) error {
	if wld.Fragments == nil {
		wld.Fragments = []model.FragmentReadWriter{}
	}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	tag.NewWithCoder(dec)
	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		return fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	wld.Version = dec.Uint32()
	tag.Mark("red", "header")

	wld.IsOldWorld = false
	switch wld.Version {
	case 0x00015500:
		wld.IsOldWorld = true
	case 0x1000C800:
		wld.IsOldWorld = false
	default:
		return fmt.Errorf("unknown wld version %d", wld.Version)
	}

	fragmentCount := dec.Uint32()
	tag.Mark("blue", "fragcount")

	wld.BspRegionCount = dec.Uint32() //bspRegionCount
	tag.Mark("green", "bspRegionCount")
	wld.Unk2 = dec.Uint32() //unk2
	tag.Mark("lime", "unk2")
	hashSize := dec.Uint32()
	tag.Mark("green", "hashsize")
	wld.Unk3 = dec.Uint32() //unk3
	tag.Mark("lime", "unk3")
	hashRaw := dec.Bytes(int(hashSize))
	nameData := helper.ReadStringHash(hashRaw)
	tag.Mark("red", "namehash")

	names = []*nameEntry{}
	chunk := []rune{}
	lastOffset := 0
	//nameBuf = []byte{}
	for i, b := range nameData {
		if b == 0 {
			names = append(names, &nameEntry{name: string(chunk), offset: lastOffset})

			//nameBuf = append(nameBuf, []byte(string(chunk))...)
			//nameBuf = append(nameBuf, 0)
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		if i == len(nameData)-1 {
			break // some times there's garbage at the end
		}
		chunk = append(chunk, b)
	}

	nameBuf = hashRaw

	fragments, err := readFragments(fragmentCount, r)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	for i := uint32(1); i <= fragmentCount; i++ {
		data := fragments[i-1]
		r := bytes.NewReader(data)

		reader := NewFrag(r)
		if reader == nil {
			return fmt.Errorf("unknown fragment at offset %d", i)
		}

		err = reader.Read(r)
		if err != nil {
			return fmt.Errorf("frag %d 0x%x (%s) read: %w", i, reader.FragCode(), FragName(int(reader.FragCode())), err)
		}
		wld.Fragments = append(wld.Fragments, reader)

	}

	return nil
}

// rawFrags is user by tests to compare for writer
func (wld *Wld) rawFrags(r io.ReadSeeker) ([][]byte, error) {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	tag.NewWithCoder(dec)
	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		return nil, fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	wld.Version = dec.Uint32()
	tag.Mark("red", "header")

	wld.IsOldWorld = false
	switch wld.Version {
	case 0x00015500:
		wld.IsOldWorld = true
	case 0x1000C800:
		wld.IsOldWorld = false
	default:
		return nil, fmt.Errorf("unknown wld version %d", wld.Version)
	}

	fragmentCount := dec.Uint32()

	wld.BspRegionCount = dec.Uint32() //bspRegionCount
	wld.Unk2 = dec.Uint32()           //unk2
	hashSize := dec.Uint32()
	wld.Unk3 = dec.Uint32() //unk3
	hashRaw := dec.Bytes(int(hashSize))
	nameData := helper.ReadStringHash(hashRaw)

	names = []*nameEntry{}
	chunk := []rune{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names = append(names, &nameEntry{name: string(chunk), offset: lastOffset})
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		if i == len(nameData)-1 {
			break // some times there's garbage at the end
		}
		chunk = append(chunk, b)
	}

	nameBuf = hashRaw

	fragments, err := readFragments(fragmentCount, r)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	return fragments, nil
}

// readFragments convert frag data to structs
func readFragments(fragmentCount uint32, r io.ReadSeeker) (fragments [][]byte, err error) {

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(fragmentCount); fragOffset++ {
		fragSize := dec.Uint32()
		totalFragSize += fragSize

		fragCode := dec.Bytes(4)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("frag position seek %d/%d: %w", fragOffset, fragmentCount, err)
		}

		data := make([]byte, fragSize)
		_, err = r.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read frag %d/%d: %w", fragOffset, fragmentCount, err)
		}

		data = append(fragCode, data...)

		fragments = append(fragments, data)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek end of frag %d/%d: %w", fragOffset, fragmentCount, err)
		}
	}

	if dec.Error() != nil {
		return nil, fmt.Errorf("read: %w", dec.Error())
	}
	return fragments, nil
}

// SetFileName sets the name of the file
func (wld *Wld) SetFileName(name string) {
	wld.MetaFileName = name
}

// FileName returns the name of the file
func (wld *Wld) FileName() string {
	return wld.MetaFileName
}
