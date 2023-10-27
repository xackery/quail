package zon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/tag"
)

// Encode writes a zon file
func Encode(zone *common.Zone, version uint32, w io.Writer) error {
	var err error
	zone.NameClear()

	tag.New()
	wEnc := encdec.NewEncoder(w, binary.LittleEndian)
	if version >= 4 {
		wEnc.String("EQTZ")
	} else {
		wEnc.String("EQGZ")
	}

	wEnc.Uint32(version)

	// rest of writer is written later
	buf := &bytes.Buffer{}
	enc := encdec.NewEncoder(buf, binary.LittleEndian)

	for _, object := range zone.Objects {
		zone.NameAdd(object.Name)
		zone.NameAdd(object.ModelName)
	}

	for i, modelName := range zone.Models {
		nameOffset := zone.NameIndex(modelName)
		if nameOffset == -1 {
			nameOffset = zone.NameAdd(modelName)
		}
		enc.Uint32(uint32(nameOffset))
		if i == 0 && strings.HasSuffix(strings.ToUpper(modelName), ".TER") {
			zone.NameAdd(strings.TrimSuffix(modelName, ".TER"))
		}
	}

	for _, object := range zone.Objects {
		nameOffset := zone.NameIndex(object.Name)
		if nameOffset == -1 {
			nameOffset = zone.NameAdd(object.Name)
		}
		enc.Uint32(uint32(nameOffset))

		enc.Float32(object.Position.X)
		enc.Float32(object.Position.Y)
		enc.Float32(object.Position.Z)

		enc.Float32(object.Rotation.X)
		enc.Float32(object.Rotation.Y)
		enc.Float32(object.Rotation.Z)

		enc.Float32(object.Scale)
	}

	for _, region := range zone.Regions {
		nameOffset := zone.NameIndex(region.Name)
		if nameOffset == -1 {
			nameOffset = zone.NameAdd(region.Name)
		}
		enc.Uint32(uint32(nameOffset))

		enc.Float32(region.Center.X)
		enc.Float32(region.Center.Y)
		enc.Float32(region.Center.Z)

		enc.Float32(region.Unknown.X)
		enc.Float32(region.Unknown.Y)
		enc.Float32(region.Unknown.Z)

		enc.Float32(region.Extent.X)
		enc.Float32(region.Extent.Y)
		enc.Float32(region.Extent.Z)
	}

	for _, light := range zone.Lights {
		nameOffset := zone.NameIndex(light.Name)
		if nameOffset == -1 {
			nameOffset = zone.NameAdd(light.Name)
		}
		enc.Uint32(uint32(nameOffset))

		enc.Float32(light.Position.X)
		enc.Float32(light.Position.Y)
		enc.Float32(light.Position.Z)

		enc.Float32(light.Color.X)
		enc.Float32(light.Color.Y)
		enc.Float32(light.Color.Z)

		enc.Float32(light.Radius)

	}

	nameData := zone.NameData()
	wEnc.Uint32(uint32(len(nameData)))
	wEnc.Uint32(uint32(len(zone.Models)))
	wEnc.Uint32(uint32(len(zone.Objects)))
	wEnc.Uint32(uint32(len(zone.Regions)))
	wEnc.Uint32(uint32(len(zone.Lights)))

	wEnc.Bytes(nameData)
	wEnc.Bytes(buf.Bytes()) // write delayed info

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	//log.Debugf("%s encoded %d objects, %d regions, %d lights", zone.Header.Name, len(zone.Objects), len(zone.Regions), len(zone.Lights))
	return nil
}
