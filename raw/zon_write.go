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
	wEnc := encdec.NewEncoder(w, binary.LittleEndian)
	if zon.Version >= 4 {
		wEnc.String("EQTZ")
	} else {
		wEnc.String("EQGZ")
	}

	wEnc.Uint32(uint32(zon.Version))

	// rest of writer is written later
	buf := &bytes.Buffer{}
	enc := encdec.NewEncoder(buf, binary.LittleEndian)

	for _, object := range zon.Objects {
		NameAdd(object.Name)
		NameAdd(object.ModelName)
	}

	for i, modelName := range zon.Models {
		nameOffset := NameIndex(modelName)
		if nameOffset == -1 {
			nameOffset = NameAdd(modelName)
		}
		enc.Uint32(uint32(nameOffset))
		if i == 0 && strings.HasSuffix(strings.ToUpper(modelName), ".TER") {
			NameAdd(strings.TrimSuffix(modelName, ".TER"))
		}
	}

	for _, object := range zon.Objects {
		nameOffset := NameIndex(object.Name)
		if nameOffset == -1 {
			nameOffset = NameAdd(object.Name)
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

	for _, region := range zon.Regions {
		nameOffset := NameIndex(region.Name)
		if nameOffset == -1 {
			nameOffset = NameAdd(region.Name)
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

	for _, light := range zon.Lights {
		nameOffset := NameIndex(light.Name)
		if nameOffset == -1 {
			nameOffset = NameAdd(light.Name)
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

	nameData := NameData()
	wEnc.Uint32(uint32(len(nameData)))
	wEnc.Uint32(uint32(len(zon.Models)))
	wEnc.Uint32(uint32(len(zon.Objects)))
	wEnc.Uint32(uint32(len(zon.Regions)))
	wEnc.Uint32(uint32(len(zon.Lights)))

	wEnc.Bytes(nameData)
	wEnc.Bytes(buf.Bytes()) // write delayed info

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	//log.Debugf("%s encoded %d objects, %d regions, %d lights", zon.Header.Name, len(zon.Objects), len(zon.Regions), len(zon.Lights))
	return nil
}
