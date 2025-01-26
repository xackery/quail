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
	if zon.name == nil {
		zon.name = &eqgName{}
	}
	if zon.Version >= 4 {
		return zon.WriteV4(w)
	}
	zon.name.clear()

	enc := encdec.NewEncoder(w, binary.LittleEndian)

	enc.String("EQGZ")

	enc.Uint32(uint32(zon.Version))

	// rest of writer is written later
	buf := &bytes.Buffer{}
	subEnc := encdec.NewEncoder(buf, binary.LittleEndian)

	for _, modelName := range zon.Models {
		zon.name.offsetByName(modelName)
	}

	for _, object := range zon.Objects {
		zon.name.offsetByName(object.InstanceName)
	}

	for _, region := range zon.Regions {
		zon.name.offsetByName(region.Name)
	}

	for _, modelName := range zon.Models {
		subEnc.Int32(zon.name.indexByName(modelName))
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
		subEnc.Int32(zon.name.indexByName(object.InstanceName))

		subEnc.Float32(object.Position[1]) //  y before x
		subEnc.Float32(object.Position[0])
		subEnc.Float32(object.Position[2])

		subEnc.Float32(object.Rotation[0])
		subEnc.Float32(object.Rotation[1])
		subEnc.Float32(object.Rotation[2])

		subEnc.Float32(object.Scale)
		if zon.Version >= 2 {
			subEnc.Uint32(uint32(len(object.Lits)))
			for _, lit := range object.Lits {
				subEnc.Uint8(lit[0])
				subEnc.Uint8(lit[1])
				subEnc.Uint8(lit[2])
				subEnc.Uint8(lit[3])
			}
		}
	}

	for _, region := range zon.Regions {
		subEnc.Int32(zon.name.indexByName(region.Name))

		subEnc.Float32(region.Center[0])
		subEnc.Float32(region.Center[1])
		subEnc.Float32(region.Center[2])

		subEnc.Float32(region.Unknown[0])
		subEnc.Float32(region.Unknown[1])
		subEnc.Float32(region.Unknown[2])

		subEnc.Float32(region.Extent[0])
		subEnc.Float32(region.Extent[1])
		subEnc.Float32(region.Extent[2])

		//subEnc.Uint32(region.Unk1)
		//subEnc.Uint32(region.Unk2)
	}

	for _, light := range zon.Lights {
		subEnc.Int32(zon.name.indexByName(light.Name))

		subEnc.Float32(light.Position[0])
		subEnc.Float32(light.Position[1])
		subEnc.Float32(light.Position[2])

		subEnc.Float32(light.Color[0])
		subEnc.Float32(light.Color[1])
		subEnc.Float32(light.Color[2])

		subEnc.Float32(light.Radius)

	}

	nameData := zon.name.data()
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
