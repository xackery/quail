package common

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
	Version       int
	IsOldWorld    bool
	Names         map[int32]string // used temporarily while decoding a wld
	fragments     [][]byte         // used temporarily while decoding a wld
	Materials     []*Material
	FragmentCount uint32
	Models        []*Model
}

// WldOpen prepares a wld file, loading fragments. This is usually then called by Decode
func WldOpen(r io.ReadSeeker) (*Wld, error) {
	wld := &Wld{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	tag.New()

	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		return nil, fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	wld.Version = int(dec.Uint32())

	wld.IsOldWorld = false
	switch wld.Version {
	case 0x00015500:
		wld.IsOldWorld = true
	case 0x1000C800:
		wld.IsOldWorld = false
	default:
		return nil, fmt.Errorf("unknown wld identifier %d", wld.Version)
	}
	wld.FragmentCount = dec.Uint32()
	_ = dec.Uint32() //unk1
	_ = dec.Uint32() //unk2
	hashSize := dec.Uint32()
	_ = dec.Uint32() //unk3
	tag.Add(tag.LastPos(), dec.Pos(), "red", "header")
	hashRaw := dec.Bytes(int(hashSize))
	nameData := helper.ReadStringHash(hashRaw)

	wld.Names = make(map[int32]string)
	chunk := []rune{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			wld.Names[int32(lastOffset)] = string(chunk)
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "green", "namedata")

	if dec.Error() != nil {
		return nil, fmt.Errorf("decode: %w", dec.Error())
	}

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(wld.FragmentCount); fragOffset++ {

		fragSize := dec.Uint32()
		totalFragSize += fragSize

		fragCode := dec.Bytes(4)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("frag position seek %d/%d: %w", fragOffset, wld.FragmentCount, err)
		}
		data := make([]byte, fragSize)
		_, err = r.Read(data)
		if err != nil {
			return nil, fmt.Errorf("read frag %d/%d: %w", fragOffset, wld.FragmentCount, err)
		}

		data = append(fragCode, data...)

		wld.fragments = append(wld.fragments, data)

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek end of frag %d/%d: %w", fragOffset, wld.FragmentCount, err)
		}
	}
	return wld, nil
}

// Fragment returns data from a specific fragment, used primarily for tests
func (wld *Wld) Fragment(fragmentIndex int) ([]byte, error) {
	if fragmentIndex < 0 || fragmentIndex >= len(wld.fragments) {
		return nil, fmt.Errorf("fragment %d out of bounds", fragmentIndex)
	}
	return wld.fragments[fragmentIndex], nil
}

func (wld *Wld) Close() {
	wld.fragments = nil
	wld.Names = nil
	wld.Materials = nil
}
