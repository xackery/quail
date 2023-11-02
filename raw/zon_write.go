package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/tag"
)

// Encode writes a zon file
func (zon *Zon) Write(w io.Writer) error {
	var err error
	NameClear()

	tag.New()
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	if zon.Version >= 4 {
		enc.String("EQTZ")
	} else {
		enc.String("EQGZ")
	}

	enc.Uint32(uint32(zon.Version))

	// rest of writer is written later
	buf := &bytes.Buffer{}
	subEnc := encdec.NewEncoder(buf, binary.LittleEndian)

	for _, object := range zon.Objects {
		NameAdd(object.Name)
		NameAdd(object.ModelName)
	}

	for i, modelName := range zon.Models {
		nameOffset := NameIndex(modelName)
		if nameOffset == -1 {
			nameOffset = NameAdd(modelName)
		}
		subEnc.Uint32(uint32(nameOffset))
		if i == 0 && strings.HasSuffix(strings.ToUpper(modelName), ".TER") {
			NameAdd(strings.TrimSuffix(modelName, ".TER"))
		}
	}

	for _, object := range zon.Objects {
		nameOffset := NameIndex(object.Name)
		if nameOffset == -1 {
			nameOffset = NameAdd(object.Name)
		}
		subEnc.Uint32(uint32(nameOffset))

		subEnc.Float32(object.Position.X)
		subEnc.Float32(object.Position.Y)
		subEnc.Float32(object.Position.Z)

		subEnc.Float32(object.Rotation.X)
		subEnc.Float32(object.Rotation.Y)
		subEnc.Float32(object.Rotation.Z)

		subEnc.Float32(object.Scale)
	}

	for _, region := range zon.Regions {
		nameOffset := NameIndex(region.Name)
		if nameOffset == -1 {
			nameOffset = NameAdd(region.Name)
		}
		subEnc.Uint32(uint32(nameOffset))

		subEnc.Float32(region.Center.X)
		subEnc.Float32(region.Center.Y)
		subEnc.Float32(region.Center.Z)

		subEnc.Float32(region.Unknown.X)
		subEnc.Float32(region.Unknown.Y)
		subEnc.Float32(region.Unknown.Z)

		subEnc.Float32(region.Extent.X)
		subEnc.Float32(region.Extent.Y)
		subEnc.Float32(region.Extent.Z)
	}

	for _, light := range zon.Lights {
		nameOffset := NameIndex(light.Name)
		if nameOffset == -1 {
			nameOffset = NameAdd(light.Name)
		}
		subEnc.Uint32(uint32(nameOffset))

		subEnc.Float32(light.Position.X)
		subEnc.Float32(light.Position.Y)
		subEnc.Float32(light.Position.Z)

		subEnc.Float32(light.Color.X)
		subEnc.Float32(light.Color.Y)
		subEnc.Float32(light.Color.Z)

		subEnc.Float32(light.Radius)

	}

	nameData := NameData()
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(len(zon.Models)))
	enc.Uint32(uint32(len(zon.Objects)))
	enc.Uint32(uint32(len(zon.Regions)))
	enc.Uint32(uint32(len(zon.Lights)))

	enc.Bytes(nameData)
	enc.Bytes(buf.Bytes()) // write delayed info

	err = subEnc.Error()
	if err != nil {
		return fmt.Errorf("subEncode: %w", err)
	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	//log.Debugf("%s encoded %d objects, %d regions, %d lights", zon.Header.Name, len(zon.Objects), len(zon.Regions), len(zon.Lights))
	return nil
}
