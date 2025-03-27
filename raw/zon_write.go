package raw

import (
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

	for _, instance := range zon.Instances {
		zon.name.add(instance.ModelTag)
		zon.name.add(instance.InstanceTag)
	}

	for _, lights := range zon.Lights {
		zon.name.add(lights.Name)
	}

	for _, area := range zon.Areas {
		zon.name.add(area.Name)
	}

	enc.Uint32(uint32(len(zon.name.data())))
	enc.Uint32(uint32(len(zon.Models)))
	enc.Uint32(uint32(len(zon.Instances)))
	enc.Uint32(uint32(len(zon.Areas)))
	enc.Uint32(uint32(len(zon.Lights)))

	enc.Bytes(zon.name.data())

	for _, modelTag := range zon.Models {
		enc.Int32(zon.name.offsetByName(modelTag))
	}

	for _, instance := range zon.Instances {
		isFound := false
		for i, name := range zon.Models {
			if name != instance.ModelTag {
				continue
			}
			enc.Int32(int32(i))
			isFound = true
			break
		}
		if !isFound {
			return fmt.Errorf("instance %s ref to model %s not found", instance.InstanceTag, instance.ModelTag)
		}
		enc.Int32(zon.name.offsetByName(instance.InstanceTag))

		enc.Float32(instance.Translation[0])
		enc.Float32(instance.Translation[1])
		enc.Float32(instance.Translation[2])

		enc.Float32(instance.Rotation[0])
		enc.Float32(instance.Rotation[1])
		enc.Float32(instance.Rotation[2])

		enc.Float32(instance.Scale)
		if zon.Version > 1 {
			enc.Uint32(uint32(len(instance.Lits)))
			for _, lit := range instance.Lits {
				enc.Uint32(lit)
			}
		}
	}

	for _, area := range zon.Areas {
		enc.Int32(zon.name.offsetByName(area.Name))

		enc.Float32(area.Center[0])
		enc.Float32(area.Center[1])
		enc.Float32(area.Center[2])

		enc.Float32(area.Orientation[0])
		enc.Float32(area.Orientation[1])
		enc.Float32(area.Orientation[2])

		enc.Float32(area.Extents[0])
		enc.Float32(area.Extents[1])
		enc.Float32(area.Extents[2])
	}

	for _, light := range zon.Lights {
		enc.Int32(zon.name.offsetByName(light.Name))

		enc.Float32(light.Position[0])
		enc.Float32(light.Position[1])
		enc.Float32(light.Position[2])

		enc.Float32(light.Color[0])
		enc.Float32(light.Color[1])
		enc.Float32(light.Color[2])

		enc.Float32(light.Radius)
	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("subEncode: %w", err)
	}

	return nil
}
