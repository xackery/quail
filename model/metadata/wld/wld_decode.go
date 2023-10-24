package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/tag"
)

var (
	isOldWorld bool
)

// Decode decodes a wld file that was prepped by Load
func Decode(wld *common.Wld, r io.ReadSeeker) error {

	if wld.Fragments == nil {
		wld.Fragments = make(map[int]common.FragmentReader)
	}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		return fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	wld.Header.Version = int(dec.Uint32())

	wld.IsOldWorld = false
	switch wld.Header.Version {
	case 0x00015500:
		wld.IsOldWorld = true
		isOldWorld = true
	case 0x1000C800:
		wld.IsOldWorld = false
	default:
		return fmt.Errorf("unknown wld version %d", wld.Header.Version)
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
	wld.SetNames(names)
	tag.Add(tag.LastPos(), dec.Pos(), "green", "namedata")

	fragments, err := loadFragments(fragmentCount, r)
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
			return fmt.Errorf("frag %d 0x%x (%s) decode: %w", i, fragCode, common.FragName(int(fragCode)), err)
		}

		//fmt.Println("frag", i, fragCode, common.FragName(int(fragCode)))
	}

	return nil
}
