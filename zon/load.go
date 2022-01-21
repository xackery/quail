package zon

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/dump"
)

func (e *ZON) Load(r io.ReadSeeker) error {
	var err error

	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'E', 'Q', 'G', 'Z'} {
		return fmt.Errorf("header does not match EQGZ")
	}

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("read header version: %w", err)
	}
	dump.Hex(version, "version=%d", version)
	if version != 1 {
		return fmt.Errorf("version is %d, wanted 1", version)
	}

	nameLength := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &nameLength)
	if err != nil {
		return fmt.Errorf("read name length: %w", err)
	}
	dump.Hex(nameLength, "nameLength=%d", nameLength)

	modelCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &modelCount)
	if err != nil {
		return fmt.Errorf("read model count: %w", err)
	}
	dump.Hex(modelCount, "modelCount=%d", modelCount)

	objectCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &objectCount)
	if err != nil {
		return fmt.Errorf("read object count: %w", err)
	}
	dump.Hex(objectCount, "objectCount=%d", objectCount)

	regionCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &regionCount)
	if err != nil {
		return fmt.Errorf("read region count: %w", err)
	}
	dump.Hex(regionCount, "regionCount=%d", regionCount)

	lightCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &lightCount)
	if err != nil {
		return fmt.Errorf("read light count: %w", err)
	}
	dump.Hex(lightCount, "lightCount=%d", lightCount)

	nameData := make([]byte, nameLength)

	err = binary.Read(r, binary.LittleEndian, &nameData)
	if err != nil {
		return fmt.Errorf("read nameData: %w", err)
	}

	type nameInfo struct {
		name   string
		offset uint32
	}
	names := []*nameInfo{}

	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names = append(names, &nameInfo{
				name:   string(chunk),
				offset: uint32(lastOffset),
			})
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}
	dump.HexRange(nameData, int(nameLength), "nameData=(%d bytes, %d entries)", nameLength, len(names))

	for i := 0; i < int(modelCount); i++ {
		modelNameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &modelNameOffset)
		if err != nil {
			return fmt.Errorf("model %d name offset: %w", i, err)
		}
		name := ""
		for _, val := range names {
			if val.offset != modelNameOffset {
				continue
			}
			name = val.name
		}
		if name == "" {
			return fmt.Errorf("model %d name at offset 0x%x not found", i, modelNameOffset)
		}

		err = e.AddModel(name)
		if err != nil {
			return fmt.Errorf("addModel %s: %w", name, err)
		}
	}

	for i := 0; i < int(objectCount); i++ {

		modelID := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &modelID)
		if err != nil {
			return fmt.Errorf("object %d modelID: %w", i, err)
		}
		if len(names) <= int(modelID) {
			return fmt.Errorf("modelID greater than names")
		}
		modelName := names[int(modelID)].name
		//dump.Hex(modelID, "modelID=%d(%s)", modelID, names[int(modelID)].name)

		objectNameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &objectNameOffset)
		if err != nil {
			return fmt.Errorf("object %d name ID: %w", i, err)
		}
		name := ""
		for _, val := range names {
			if val.offset != objectNameOffset {
				continue
			}
			name = val.name
		}
		if name == "" {
			return fmt.Errorf("model %d name at offset 0x%x not found", i, objectNameOffset)
		}
		//dump.Hex(objectNameOffset, "objectNameOffset=0x%x(%s)", objectNameOffset, name)

		pos := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("object %d pos: %w", i, err)
		}
		//dump.Hex(pos, "pos=%+v", pos)

		rot := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &rot)
		if err != nil {
			return fmt.Errorf("object %d rot: %w", i, err)
		}
		//dump.Hex(rot, "rot=%+v", rot)

		scale := float32(0)
		err = binary.Read(r, binary.LittleEndian, &scale)
		if err != nil {
			return fmt.Errorf("object %d scale: %w", i, err)
		}
		//dump.Hex(scale, "scale=%0.2f", scale)

		err = e.AddObject(modelName, name, pos, rot, scale)
		if err != nil {
			return fmt.Errorf("addObject %s: %w", name, err)
		}
	}
	dump.HexRange([]byte{1, 2}, int(objectCount*36), "objectChunk=(%d bytes, %d entries)", int(objectCount*36), objectCount)

	for i := 0; i < int(regionCount); i++ {
		regionNameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &regionNameOffset)
		if err != nil {
			return fmt.Errorf("region %d name ID: %w", i, err)
		}
		name := ""
		for _, val := range names {
			if val.offset != regionNameOffset {
				continue
			}
			name = val.name
		}
		if name == "" {
			return fmt.Errorf("model %d name at offset 0x%x not found", i, regionNameOffset)
		}
		//dump.Hex(regionNameOffset, "regionNameOffset=0x%x(%s)", regionNameOffset, name)

		center := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &center)
		if err != nil {
			return fmt.Errorf("region %d center: %w", i, err)
		}
		//dump.Hex(center, "center=%+v", center)

		unknown := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &unknown)
		if err != nil {
			return fmt.Errorf("region %d unknown: %w", i, err)
		}
		//dump.Hex(unknown, "unknown=%+v", unknown)

		extent := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &extent)
		if err != nil {
			return fmt.Errorf("region %d extent: %w", i, err)
		}
		//dump.Hex(extent, "extent=%+v", extent)
		err = e.AddRegion(name, center, unknown, extent)
		if err != nil {
			return fmt.Errorf("addRegion %s: %w", name, err)
		}
	}
	dump.HexRange([]byte{1, 2}, int(regionCount*40), "regionChunk=(%d bytes, %d entries)", int(regionCount*36), regionCount)

	for i := 0; i < int(lightCount); i++ {
		lightNameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &lightNameOffset)
		if err != nil {
			return fmt.Errorf("light %d name ID: %w", i, err)
		}
		name := ""
		for _, val := range names {
			if val.offset != lightNameOffset {
				continue
			}
			name = val.name
		}
		if name == "" {
			return fmt.Errorf("model %d name at offset 0x%x not found", i, lightNameOffset)
		}
		//dump.Hex(lightNameOffset, "lightNameOffset=0x%x(%s)", lightNameOffset, name)

		pos := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("light %d pos: %w", i, err)
		}
		//dump.Hex(pos, "pos=%+v", pos)

		color := math32.Color{}
		err = binary.Read(r, binary.LittleEndian, &color)
		if err != nil {
			return fmt.Errorf("light %d color: %w", i, err)
		}
		//dump.Hex(color, "color=%+v", color)

		radius := float32(0)
		err = binary.Read(r, binary.LittleEndian, &radius)
		if err != nil {
			return fmt.Errorf("light %d radius: %w", i, err)
		}
		//dump.Hex(radius, "radius=%+v", radius)

		err = e.AddLight(name, pos, color, radius)
		if err != nil {
			return fmt.Errorf("addLight %s: %w", name, err)
		}

	}
	dump.HexRange([]byte{1, 2}, int(lightCount*32), "lightChunk=(%d bytes, %d entries)", int(lightCount*32), lightCount)

	return nil
}
