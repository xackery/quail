package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/tag"
)

type Wld struct {
	Version    uint32                 `yaml:"version"`
	IsOldWorld bool                   `yaml:"is_old_world"`
	Fragments  map[int]FragmentReader `yaml:"fragments,omitempty"`
}

// Read decodes a wld file that was prepped by Load
func (wld *Wld) Read(r io.ReadSeeker) error {
	if wld.Fragments == nil {
		wld.Fragments = make(map[int]FragmentReader)
	}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		return fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	wld.Version = dec.Uint32()

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
	_ = dec.Uint32() //unk1
	_ = dec.Uint32() //unk2
	hashSize := dec.Uint32()
	_ = dec.Uint32() //unk3
	tag.Add(tag.LastPos(), dec.Pos(), "red", "header")
	hashRaw := dec.Bytes(int(hashSize))
	nameData := helper.ReadStringHash(hashRaw)

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

	NamesSet(names)
	tag.Add(tag.LastPos(), dec.Pos(), "green", "namedata")

	fragments, err := readFragments(fragmentCount, r)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	tag.New()
	for i := uint32(1); i <= fragmentCount; i++ {
		data := fragments[i-1]
		r := bytes.NewReader(data)

		dec := encdec.NewDecoder(r, binary.LittleEndian)

		fragCode := dec.Int32()

		decoder, ok := decoders[int(fragCode)]
		if !ok {
			return fmt.Errorf("frag %d 0x%x decode: unsupported fragment", i, fragCode)
		}

		wld.Fragments[int(i)], err = decoder(r)
		if err != nil {
			return fmt.Errorf("frag %d 0x%x (%s) decode: %w", i, fragCode, FragName(int(fragCode)), err)
		}

		//fmt.Println("frag", i, fragCode, FragName(int(fragCode)))
	}

	return nil
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
		return nil, fmt.Errorf("decode: %w", dec.Error())
	}
	return fragments, nil
}
