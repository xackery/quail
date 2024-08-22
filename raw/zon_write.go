package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Encode writes a zon file
func (zon *Zon) Write(w io.Writer) error {
	var err error
	if zon.Version >= 4 {
		return zon.WriteV4(w)
	}
	zon.NameClear()

	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.String("EQGZ")

	enc.Uint32(uint32(zon.Version))

	// rest of writer is written later
	buf := &bytes.Buffer{}
	subEnc := encdec.NewEncoder(buf, binary.LittleEndian)

	for _, modelName := range zon.Models {
		zon.NameAdd(modelName)
	}

	for _, object := range zon.Objects {
		zon.NameAdd(object.InstanceName)
	}

	for _, region := range zon.Regions {
		zon.NameAdd(region.Name)
	}

	for _, modelName := range zon.Models {
		subEnc.Int32(zon.NameIndex(modelName))
	}

	for _, object := range zon.Objects {
		isFound := false
		for i, name := range zon.Models {
			if name != object.ModelName {
				continue
			}
			subEnc.Int32(int32(i))
			isFound = true
			break
		}
		if !isFound {
			return fmt.Errorf("object %s ref to model %s not found", object.InstanceName, object.ModelName)
		}
		subEnc.Int32(zon.NameIndex(object.InstanceName))

		subEnc.Float32(object.Position.Y) //  y before x
		subEnc.Float32(object.Position.X)
		subEnc.Float32(object.Position.Z)

		subEnc.Float32(object.Rotation.X)
		subEnc.Float32(object.Rotation.Y)
		subEnc.Float32(object.Rotation.Z)

		subEnc.Float32(object.Scale)
		if zon.Version >= 2 {
			subEnc.Uint32(uint32(len(object.Lits)))
			for _, lit := range object.Lits {
				subEnc.Uint8(lit.R)
				subEnc.Uint8(lit.G)
				subEnc.Uint8(lit.B)
				subEnc.Uint8(lit.A)
			}
		}
	}

	for _, region := range zon.Regions {
		subEnc.Int32(zon.NameIndex(region.Name))

		subEnc.Float32(region.Center.X)
		subEnc.Float32(region.Center.Y)
		subEnc.Float32(region.Center.Z)

		subEnc.Float32(region.Unknown.X)
		subEnc.Float32(region.Unknown.Y)
		subEnc.Float32(region.Unknown.Z)

		subEnc.Float32(region.Extent.X)
		subEnc.Float32(region.Extent.Y)
		subEnc.Float32(region.Extent.Z)

		//subEnc.Uint32(region.Unk1)
		//subEnc.Uint32(region.Unk2)
	}

	for _, light := range zon.Lights {
		subEnc.Int32(zon.NameIndex(light.Name))

		subEnc.Float32(light.Position.X)
		subEnc.Float32(light.Position.Y)
		subEnc.Float32(light.Position.Z)

		subEnc.Float32(light.Color.X)
		subEnc.Float32(light.Color.Y)
		subEnc.Float32(light.Color.Z)

		subEnc.Float32(light.Radius)

	}

	nameData := zon.NameData()
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

	return nil
}
