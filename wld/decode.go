package wld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/wld/fragment"
)

// Decode loads a wld file
func (e *WLD) Decode(r io.ReadSeeker) error {
	var err error
	if e == nil {
		return fmt.Errorf("wld nil")
	}
	var value uint32

	var header [4]byte
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	validHeader := [4]byte{0x02, 0x3D, 0x50, 0x54}
	if header != validHeader {
		return fmt.Errorf("header wanted 0x%x, got 0x%x", validHeader, header)
	}
	dump.Hex(header, "header=0x%x", header)

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("read identifier: %w", err)
	}

	isOldWorld := false
	switch version {
	case 0x00015500:
		isOldWorld = true
	case 0x1000C800:
		isOldWorld = false
	default:
		return fmt.Errorf("unknown wld identifier %d", value)
	}
	dump.Hex(value, "identifier=(isOld:%t)", isOldWorld)

	fragmentCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &fragmentCount)
	if err != nil {
		return fmt.Errorf("read fragmentCount: %w", err)
	}
	dump.Hex(fragmentCount, "fragmentCount=%d", fragmentCount)

	unk1 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk1)
	if err != nil {
		return fmt.Errorf("read unk1: %w", err)
	}
	dump.Hex(unk1, "unk1=%d", unk1)

	unk2 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk2)
	if err != nil {
		return fmt.Errorf("read unk2: %w", err)
	}
	dump.Hex(unk2, "unk2=%d", unk2)

	var hashSize uint32
	err = binary.Read(r, binary.LittleEndian, &hashSize)
	if err != nil {
		return fmt.Errorf("read hash size: %w", err)
	}
	dump.Hex(hashSize, "hashSize=%d", hashSize)

	unk3 := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &unk3)
	if err != nil {
		return fmt.Errorf("read unk3: %w", err)
	}
	dump.Hex(unk3, "unk3=%d", unk3)

	hashRaw := make([]byte, hashSize)

	err = binary.Read(r, binary.LittleEndian, &hashRaw)
	if err != nil {
		return fmt.Errorf("read nameData: %w", err)
	}

	nameData := decodeStringHash(hashRaw)

	names := make(map[uint32]string)

	chunk := []rune{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			chunk = []rune{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	dump.HexRange(hashRaw, int(hashSize), "nameData=(%d bytes, %d names)", hashSize, len(names))

	totalFragSize := uint32(0)
	for i := 0; i < int(fragmentCount); i++ {
		var fragSize uint32
		var fragIndex int32

		name := names[uint32(i)]
		err = binary.Read(r, binary.LittleEndian, &fragSize)
		if err != nil {
			return fmt.Errorf("read fragment size %d/%d: %w", i, fragmentCount, err)
		}
		totalFragSize += fragSize
		//dump.Hex(fragSize, "%d(%s)fragSize=%d", i, name, fragSize)
		err = binary.Read(r, binary.LittleEndian, &fragIndex)
		if err != nil {
			return fmt.Errorf("read fragment index %d/%d: %w", i, fragmentCount, err)
		}
		//dump.Hex(fragSize, "%dfragIndex=%d", i, fragIndex)

		fragPosition, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("frag position seek %d/%d: %w", i, fragmentCount, err)
		}

		buf := make([]byte, fragSize)
		_, err = r.Read(buf)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		frag, err := fragment.New(fragIndex, bytes.NewReader(buf))
		if err != nil {
			fmt.Printf("warning: fragment decode 0x%x (%d): %s\n", fragIndex, fragIndex, err)
			//return fmt.Errorf("fragment load: %w", err)
		} else {
			e.fragments = append(e.fragments, &fragmentInfo{name: name, data: frag})
		}

		_, err = r.Seek(fragPosition+int64(fragSize), io.SeekStart)
		if err != nil {
			return fmt.Errorf("seek end of frag %d/%d: %w", i, fragmentCount, err)
		}
	}
	//dump.HexRange([]byte{byte(i), byte(i) + 1}, int(fragSize), "%dfrag=%s", i, frag.FragmentType())
	dump.HexRange([]byte{0, 1}, int(totalFragSize), "fragChunk=(%d bytes, %d entries)", int(totalFragSize), len(e.fragments))
	err = e.convertFragments()
	if err != nil {
		return fmt.Errorf("convertFragments: %w", err)
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

func (e *WLD) convertFragments() error {
	type mesher interface {
		Indices() [][3]float32
		Normals() [][3]float32
		Vertices() [][3]float32
		Uvs() [][2]float32
	}

	type materialer interface {
		Name() string
		ShaderType() int
		MaterialType() int
	}

	for _, frag := range e.fragments {
		material, ok := frag.data.(materialer)
		if !ok {
			continue
		}
		err := e.MaterialAdd(material.Name(), fmt.Sprintf("%d", material.ShaderType()))
		if err != nil {
			return fmt.Errorf("materialadd: %w", err)
		}
	}

	for _, frag := range e.fragments {
		mesh, ok := frag.data.(mesher)
		if !ok {
			continue
		}
		fmt.Println(mesh)

		/*for _, index := range mesh.Indices() {
			name, err := e.MaterialByID(index.)
			e.faces = append(e.faces, &common.Face{
				Index: [3]uint32{uint32(index.X), uint32(index.Y), uint32(index.Z)},
			})
		}*/
	}
	return nil
}
