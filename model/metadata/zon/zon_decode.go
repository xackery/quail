package zon

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/log"
)

// Decode decodes a ZON file
func Decode(zone *common.Zone, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	// Read header
	header := dec.StringFixed(4)
	dump.Hex(header, "header=%s", header)
	if header != "EQGZ" && header != "EQTZ" {
		return fmt.Errorf("header %s does not match EQGZ or EQTZ", header)
	}

	if header == "EQTZ" {
		return fmt.Errorf("type 4 zone not supported")
		//return e.eqgtzpDecode(r)
	}

	version := dec.Uint32()
	zone.Version = int(version)
	if version != 1 {
		return fmt.Errorf("version is %d, wanted 1", version)
	}

	nameLength := dec.Uint32()
	modelCount := dec.Uint32()
	objectCount := dec.Uint32()
	regionCount := dec.Uint32()
	lightCount := dec.Uint32()

	nameData := dec.Bytes(int(nameLength))
	names := make(map[uint32]string)
	namesIndexed := []string{}

	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			namesIndexed = append(namesIndexed, string(chunk))
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}

	var ok bool
	for i := 0; i < int(modelCount); i++ {
		nameOffset := dec.Uint32()
		name, ok := names[nameOffset]
		if !ok {
			return fmt.Errorf("model nameOffset %d not found", nameOffset)
		}

		zone.Models = append(zone.Models, name)
	}

	for i := 0; i < int(objectCount); i++ {
		object := common.Object{}
		nameIndex := dec.Uint32()

		if nameIndex >= uint32(len(namesIndexed)) {
			return fmt.Errorf("object nameIndex %d out of range", nameIndex)
		}

		object.Name = namesIndexed[nameIndex]

		nameOffset := dec.Uint32()
		object.ModelName, ok = names[nameOffset]
		if !ok {
			return fmt.Errorf("object modelNameOffset %d not found", nameOffset)
		}

		object.Position.X = dec.Float32()
		object.Position.Y = dec.Float32()
		object.Position.Z = dec.Float32()

		object.Rotation.X = dec.Float32()
		object.Rotation.Y = dec.Float32()
		object.Rotation.Z = dec.Float32()

		object.Scale = dec.Float32()
		zone.Objects = append(zone.Objects, object)
	}

	for i := 0; i < int(regionCount); i++ {
		region := common.Region{}

		nameOffset := dec.Uint32()
		region.Name, ok = names[nameOffset]
		if !ok {
			return fmt.Errorf("region nameOffset %d not found", nameOffset)
		}

		region.Center.X = dec.Float32()
		region.Center.Y = dec.Float32()
		region.Center.Z = dec.Float32()

		region.Unknown.X = dec.Float32()
		region.Unknown.Y = dec.Float32()
		region.Unknown.Z = dec.Float32()

		region.Extent.X = dec.Float32()
		region.Extent.Y = dec.Float32()
		region.Extent.Z = dec.Float32()

		zone.Regions = append(zone.Regions, region)
	}

	for i := 0; i < int(lightCount); i++ {
		light := common.Light{}

		nameOffset := dec.Uint32()
		light.Name, ok = names[nameOffset]
		if !ok {
			return fmt.Errorf("light nameOffset %d not found", nameOffset)
		}

		light.Position.X = dec.Float32()
		light.Position.Y = dec.Float32()
		light.Position.Z = dec.Float32()

		light.Color.X = dec.Float32()
		light.Color.Y = dec.Float32()
		light.Color.Z = dec.Float32()

		light.Radius = dec.Float32()

		zone.Lights = append(zone.Lights, light)
	}

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	log.Debugf("%s (zon) decoded %d model refs, %d object refs, %d regions, %d lights", zone.Name, len(zone.Models), len(zone.Objects), len(zone.Regions), len(zone.Lights))
	return nil
}
