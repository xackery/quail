package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/log"
)

type encoderdecoder struct {
	name   string
	decode func(r io.ReadSeeker, fragOffset int) error
	encode func(w io.Writer, fragOffset int) error
}

// Decode loads a wld file
func (e *WLD) Decode(r io.ReadSeeker) error {

	var err error
	if e == nil {
		return fmt.Errorf("wld nil")
	}

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	fragmentCount, err := e.readHeader(r)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	totalFragSize := uint32(0)
	for fragOffset := 0; fragOffset < int(fragmentCount); fragOffset++ {
		fragSize := dec.Uint32()
		totalFragSize += fragSize
		//dump.Hex(fragSize, "%d(%s)fragSize=%d", i, name, fragSize)
		fragCode := dec.Int32()
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

		parser, ok := e.packs[fragCode]
		if !ok {
			log.Warnf("warning: unknown fragCode %d at offset %d", fragCode, fragOffset)
		} else {
			err = parser.decode(bytes.NewReader(buf), fragOffset)
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

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

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

func encodeStringHash(hash string) []byte {
	hashKey := []byte{0x95, 0x3A, 0xC5, 0x2A, 0x95, 0x7A, 0x95, 0x6A}
	out := make([]byte, len(hash))
	for i := 0; i < len(hash); i++ {
		out[i] = hash[i] ^ hashKey[i%8]
	}
	return out
}

func (e *WLD) readHeader(r io.ReadSeeker) (fragmentCount uint32, err error) {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.Bytes(4)
	validHeader := []byte{0x02, 0x3D, 0x50, 0x54}
	if !bytes.Equal(header, validHeader) {
		err = fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
		return
	}
	dump.Hex(header, "header=0x%x", header)

	e.version = dec.Uint32()

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
	fragmentCount = dec.Uint32()
	dump.Hex(fragmentCount, "fragmentCount=%d", fragmentCount)
	unk1 := dec.Uint32()
	dump.Hex(unk1, "unk1=%d", unk1)
	unk2 := dec.Uint32()
	dump.Hex(unk2, "unk2=%d", unk2)
	hashSize := dec.Uint32()
	dump.Hex(hashSize, "hashSize=%d", hashSize)
	unk3 := dec.Uint32()
	dump.Hex(unk3, "unk3=%d", unk3)
	hashRaw := dec.Bytes(int(hashSize))
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

	if dec.Error() != nil {
		err = fmt.Errorf("read header: %w", dec.Error())
		return
	}

	dump.HexRange(hashRaw, int(hashSize), "nameData=(%d bytes, %d names)", hashSize, len(names))
	return
}
