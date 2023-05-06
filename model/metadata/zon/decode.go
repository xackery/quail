package zon

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// Decode decodes a zon file
func (e *ZON) Decode(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	// Read header
	header := dec.StringFixed(4)
	dump.Hex(header, "header=%s", header)
	if header != "EQGZ" && header != "EQTZ" {
		return fmt.Errorf("header %s does not match EQGZ or EQTZ", header)
	}

	if header == "EQTZ" {
		return e.eqgtzpDecode(r)
	}

	e.version = dec.Uint32()
	dump.Hex(e.version, "version=%d", e.version)
	if e.version != 1 {
		return fmt.Errorf("version is %d, wanted 1", e.version)
	}

	nameLength := dec.Uint32()
	dump.Hex(nameLength, "nameLength=%d", nameLength)
	modelCount := dec.Uint32()
	dump.Hex(modelCount, "modelCount=%d", modelCount)
	objectCount := dec.Uint32()
	dump.Hex(objectCount, "objectCount=%d", objectCount)
	regionCount := dec.Uint32()
	dump.Hex(regionCount, "regionCount=%d", regionCount)
	lightCount := dec.Uint32()
	dump.Hex(lightCount, "lightCount=%d", lightCount)

	nameData := dec.Bytes(int(nameLength))
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
		modelNameOffset := dec.Uint32()
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

		e.models = append(e.models, model{name: name})
	}
	dump.HexRange([]byte{1, 2}, int(modelCount*4), "modelChunk=(%d bytes, %d entries)", int(modelCount*4), objectCount)

	for i := 0; i < int(objectCount); i++ {
		modelID := dec.Uint32()
		if len(names) <= int(modelID) {
			return fmt.Errorf("modelID 0x%x greater than names", modelID)
		}
		modelName := names[int(modelID)].name
		dump.Hex(modelID, "%dmodelID=%d(%s)", i, modelID, names[int(modelID)].name)
		modelName = strings.ToLower(modelName)
		objectNameOffset := dec.Uint32()
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
		position.X = dec.Float32()
		position.Y = dec.Float32()
		position.Z = dec.Float32()
		dump.Hex(position, "%dposition=%+v", i, position)
		rotation := geo.Vector3{}
		rotation.X = dec.Float32()
		rotation.Y = dec.Float32()
		rotation.Z = dec.Float32()
		dump.Hex(rotation, "%drotation=%+v", i, rotation)
		scale := dec.Float32()
		dump.Hex(scale, "scale=%0.8f", scale)

		e.objectManager.Add(geo.Object{
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
		region := region{}

		regionNameOffset := dec.Uint32()
		for _, val := range names {
			if val.offset != regionNameOffset {
				continue
			}
			region.name = val.name
		}
		if region.name == "" {
			return fmt.Errorf("model %d name at offset 0x%x not found", i, regionNameOffset)
		}
		//dump.Hex(regionNameOffset, "regionNameOffset=0x%x(%s)", regionNameOffset, name)

		region.center.X = dec.Float32()
		region.center.Y = dec.Float32()
		region.center.Z = dec.Float32()
		//dump.Hex(center, "center=%+v", center)
		region.unknown.X = dec.Float32()
		region.unknown.Y = dec.Float32()
		region.unknown.Z = dec.Float32()
		//dump.Hex(unknown, "unknown=%+v", unknown)
		region.extent.X = dec.Float32()
		region.extent.Y = dec.Float32()
		region.extent.Z = dec.Float32()
		//dump.Hex(extent, "extent=%+v", extent)

		e.regions = append(e.regions, region)
	}
	dump.HexRange([]byte{1, 2}, int(regionCount*40), "regionChunk=(%d bytes, %d entries)", int(regionCount*36), regionCount)

	for i := 0; i < int(lightCount); i++ {
		lit := light{}
		lightNameOffset := dec.Uint32()
		for _, val := range names {
			if val.offset != lightNameOffset {
				continue
			}
			lit.name = val.name
		}
		if lit.name == "" {
			return fmt.Errorf("model %d name at offset 0x%x not found", i, lightNameOffset)
		}
		//dump.Hex(lightNameOffset, "lightNameOffset=0x%x(%s)", lightNameOffset, name)
		lit.position.X = dec.Float32()
		lit.position.Y = dec.Float32()
		lit.position.Z = dec.Float32()
		//dump.Hex(position, "position=%+v", position)
		lit.color.X = dec.Float32()
		lit.color.Y = dec.Float32()
		lit.color.Z = dec.Float32()
		//dump.Hex(color, "color=%+v", color)
		lit.radius = dec.Float32()
		//dump.Hex(radius, "radius=%+v", radius)

		e.lights = append(e.lights, lit)

	}
	dump.HexRange([]byte{1, 2}, int(lightCount*32), "lightChunk=(%d bytes, %d entries)", int(lightCount*32), lightCount)

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	log.Debugf("%s is version %d and has %d objects, %d regions, %d lights", e.name, e.version, len(e.objectManager.Objects()), len(e.regions), len(e.lights))
	return nil
}
