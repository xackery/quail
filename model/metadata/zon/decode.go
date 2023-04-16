package zon

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/model/geo"
)

func (e *ZON) Decode(r io.ReadSeeker) error {
	var err error

	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'E', 'Q', 'G', 'Z'} && header != [4]byte{'E', 'Q', 'T', 'Z'} {
		return fmt.Errorf("header does not match EQGZ/EQTZ")
	}

	if header == [4]byte{'E', 'Q', 'T', 'Z'} {
		return e.eqgtzpDecode(r)
	}

	err = binary.Read(r, binary.LittleEndian, &e.version)
	if err != nil {
		return fmt.Errorf("read header version: %w", err)
	}
	dump.Hex(e.version, "version=%d", e.version)
	if e.version != 1 {
		return fmt.Errorf("version is %d, wanted 1", e.version)
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
		name = strings.ToLower(name)

		e.models = append(e.models, &model{name: name})
	}
	dump.HexRange([]byte{1, 2}, int(modelCount*4), "modelChunk=(%d bytes, %d entries)", int(modelCount*4), objectCount)

	for i := 0; i < int(objectCount); i++ {

		modelID := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &modelID)
		if err != nil {
			return fmt.Errorf("object %d modelID: %w", i, err)
		}
		if len(names) <= int(modelID) {
			return fmt.Errorf("modelID 0x%x greater than names", modelID)
		}
		modelName := names[int(modelID)].name
		dump.Hex(modelID, "%dmodelID=%d(%s)", i, modelID, names[int(modelID)].name)
		modelName = strings.ToLower(modelName)
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
		dump.Hex(objectNameOffset, "%dobjectNameOffset=0x%x(%s)", i, objectNameOffset, name)
		name = strings.ToLower(name)
		position := geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &position)
		if err != nil {
			return fmt.Errorf("object %d position: %w", i, err)
		}
		dump.Hex(position, "%dposition=%+v", i, position)

		rotation := geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &rotation)
		if err != nil {
			return fmt.Errorf("object %d rotation: %w", i, err)
		}
		dump.Hex(rotation, "rotation=%+v", rotation)

		scale := float32(0)
		err = binary.Read(r, binary.LittleEndian, &scale)
		if err != nil {
			return fmt.Errorf("object %d scale: %w", i, err)
		}
		dump.Hex(scale, "scale=%0.3f", scale)

		e.objectManager.Add(&geo.Object{
			Name:      name,
			ModelName: modelName,
			Position:  position,
			Rotation:  rotation,
			Scale:     scale,
			FileType:  "",
			FileName:  "",
		})

	}
	//dump.HexRange([]byte{1, 2}, int(objectCount*36), "objectChunk=(%d bytes, %d entries)", int(objectCount*36), objectCount)

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

		center := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, center)
		if err != nil {
			return fmt.Errorf("region %d center: %w", i, err)
		}
		//dump.Hex(center, "center=%+v", center)

		unknown := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, unknown)
		if err != nil {
			return fmt.Errorf("region %d unknown: %w", i, err)
		}
		//dump.Hex(unknown, "unknown=%+v", unknown)

		extent := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, extent)
		if err != nil {
			return fmt.Errorf("region %d extent: %w", i, err)
		}
		//dump.Hex(extent, "extent=%+v", extent)

		e.regions = append(e.regions, &region{name: name, center: center, unknown: unknown, extent: extent})

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

		pos := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, pos)
		if err != nil {
			return fmt.Errorf("light %d pos: %w", i, err)
		}
		//dump.Hex(pos, "pos=%+v", pos)

		color := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, color)
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

		e.lights = append(e.lights, &light{name: name, position: pos, color: color, radius: radius})

	}
	dump.HexRange([]byte{1, 2}, int(lightCount*32), "lightChunk=(%d bytes, %d entries)", int(lightCount*32), lightCount)

	return nil
}

func (e *ZON) eqgtzpDecode(r io.ReadSeeker) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, data[0 : i+1], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "*") {
			continue
		}
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		switch parts[0] {
		case "NAME":
			e.models = append(e.models, &model{name: parts[1], baseName: parts[1]})
		case "VERSION":
			e.version = helper.AtoU32(parts[1])
			if e.version == 0 {
				return fmt.Errorf("invalid version on eqtzp: %s", parts[1])
			}
		}
	}

	return nil
}
