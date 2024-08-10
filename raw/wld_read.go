package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model"
	"github.com/xackery/quail/raw/rawfrag"
	"github.com/xackery/quail/tag"
)

type Wld struct {
	MetaFileName string                     `yaml:"file_name"`
	Version      uint32                     `yaml:"version"`
	IsNewWorld   bool                       `yaml:"is_old_world"`
	Fragments    []model.FragmentReadWriter `yaml:"fragments,omitempty"`
	Unk2         uint32                     `yaml:"unk2"`
	Unk3         uint32                     `yaml:"unk3"`
}

func (wld *Wld) Identity() string {
	return "wld"
}

// Read reads a wld file that was prepped by Load
func (wld *Wld) Read(r io.ReadSeeker) error {
	if wld.Fragments == nil {
		wld.Fragments = []model.FragmentReadWriter{&rawfrag.WldFragDefault{}}
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

	wld.IsNewWorld = false
	switch wld.Version {
	case 0x00015500:
		wld.IsNewWorld = false
	case 0x1000C800:
		wld.IsNewWorld = true
	default:
		return fmt.Errorf("unknown wld version %d", wld.Version)
	}

	fragmentCount := dec.Uint32()
	tag.Mark("blue", "fragcount")

	bspRegionCount := dec.Uint32() //bspRegionCount
	tag.Mark("green", "totalRegionCount")
	maxFragSize := dec.Uint32() // max fragment size
	hashSize := dec.Uint32()
	tag.Mark("green", "hashsize")
	stringCount := dec.Uint32() // string count
	tag.Mark("lime", "string count")
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

	if len(names) != int(stringCount)+1 {
		fmt.Printf("name count mismatch, wanted %d, got %d (ignoring, openzone?)\n", stringCount, len(names))
		//return fmt.Errorf("name count mismatch, wanted %d, got %d", stringCount, len(names))
	}

	nameBuf = hashRaw

	fragments, err := readFragments(fragmentCount, r)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	totalRegions := 0
	for i := uint32(0); i < fragmentCount; i++ {
		data := fragments[i]
		if len(data) > int(maxFragSize+4) {
			return fmt.Errorf("fragment %d (size: %d) exceeds max size %d", i, len(data), maxFragSize)
		}
		r := bytes.NewReader(data)

		reader := NewFrag(r)
		if reader == nil {
			return fmt.Errorf("unknown fragment at offset %d", i)
		}

		err = reader.Read(r, wld.IsNewWorld)
		if err != nil {
			return fmt.Errorf("frag %d 0x%x (%s) read: %w", i, reader.FragCode(), FragName(int(reader.FragCode())), err)
		}
		wld.Fragments = append(wld.Fragments, reader)

		pos, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("fragment %d (size: %d) seek: %w", i, len(data), err)
		}
		if pos != int64(len(data)) {
			isNonZero := false
			for i, bdata := range data {
				if int64(i) > pos {
					continue
				}
				if bdata != 0 {
					isNonZero = true
					break
				}
			}
			if !isNonZero {
				fmt.Printf("fragment %d seek mismatch (%d/%d) (%T)\n", i, pos, len(data), reader)

				fmt.Printf("fragment %d data: %x\n", i, data[pos:])
			}
		}
		_, ok := reader.(*rawfrag.WldFragRegion)
		if ok {
			totalRegions++
		}

	}

	if totalRegions != int(bspRegionCount) {
		return fmt.Errorf("region count mismatch, wanted %d, got %d", bspRegionCount, totalRegions)
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

	wld.IsNewWorld = true
	switch wld.Version {
	case 0x00015500:
		wld.IsNewWorld = false
	case 0x1000C800:
		wld.IsNewWorld = true
	default:
		return nil, fmt.Errorf("unknown wld version %d", wld.Version)
	}

	fragmentCount := dec.Uint32()

	_ = dec.Uint32() //bspRegionCount
	_ = dec.Uint32() // max_fragment_size
	hashSize := dec.Uint32()
	_ = dec.Uint32() // string count
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
