package zon

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/tag"
)

// Encode writes a zon file
func Encode(zone *common.Zone, version uint32, w io.Writer) error {

	modelNames := []string{}

	for _, object := range zone.Objects {
		modelNames = append(modelNames, object.Name)
	}

	names, nameData, err := zone.NameBuild(modelNames)
	if err != nil {
		return fmt.Errorf("nameBuild: %w", err)
	}

	tag.New()
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	if version >= 4 {
		enc.String("EQTZ")
	} else {
		enc.String("EQGZ")
	}

	enc.Uint32(version)
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(len(zone.Models)))
	enc.Uint32(uint32(len(zone.Objects)))
	enc.Uint32(uint32(len(zone.Regions)))
	enc.Uint32(uint32(len(zone.Lights)))

	enc.Bytes(nameData)

	for _, model := range zone.Models {
		nameOffset := -1
		for key, offset := range names {
			if key == model {
				nameOffset = int(offset)
				break
			}
		}
		//if nameOffset == -1 {
		//log.Debugf("material %s not found ignoring", o.Name)
		//}
		enc.Uint32(uint32(nameOffset))
	}

	for i, object := range zone.Objects {
		enc.Uint32(uint32(i + 1))
		nameOffset := -1
		for key, offset := range names {
			if key == object.Name {
				nameOffset = int(offset)
				break
			}
		}
		//if nameOffset == -1 {
		//log.Debugf("material %s not found ignoring", o.Name)
		//}
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
		nameOffset := -1
		for key, offset := range names {
			if key == region.Name {
				nameOffset = int(offset)
				break
			}
		}
		//if nameOffset == -1 {
		//log.Debugf("material %s not found ignoring", o.Name)
		//}
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
		nameOffset := -1
		for key, offset := range names {
			if key == light.Name {
				nameOffset = int(offset)
				break
			}
		}
		//if nameOffset == -1 {
		//log.Debugf("material %s not found ignoring", o.Name)
		//}
		enc.Uint32(uint32(nameOffset))

		enc.Float32(light.Position.X)
		enc.Float32(light.Position.Y)
		enc.Float32(light.Position.Z)

		enc.Float32(light.Color.X)
		enc.Float32(light.Color.Y)
		enc.Float32(light.Color.Z)

		enc.Float32(light.Radius)

	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	//log.Debugf("%s encoded %d objects, %d regions, %d lights", zone.Header.Name, len(zone.Objects), len(zone.Regions), len(zone.Lights))
	return nil
}
