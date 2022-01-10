package zon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/helper"
)

// Save writes a zon file to location
func (e *ZON) Save(w io.Writer) error {
	var err error

	names := []string{}
	nameBuf := bytes.NewBuffer(nil)
	dataBuf := bytes.NewBuffer(nil)

	for _, o := range e.models {
		names = append(names, o.name)
		err = helper.WriteString(nameBuf, o.name)
		if err != nil {
			return fmt.Errorf("write model nameBuf %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, uint32(len(nameBuf.Bytes())))
		if err != nil {
			return fmt.Errorf("write model dataBuf %s: %w", o.name, err)
		}
	}

	for _, o := range e.objects {
		names = append(names, o.name)
		err = helper.WriteString(nameBuf, o.name)
		if err != nil {
			return fmt.Errorf("write object nameBuf %s: %w", o.name, err)
		}

		modelID := -1
		for i, m := range e.models {
			if m.name == o.modelName {
				modelID = i
				break
			}
		}
		if modelID == -1 {
			return fmt.Errorf("object %s refers to model %s which does not exist", o.name, o.modelName)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint16(modelID))
		if err != nil {
			return fmt.Errorf("write object id %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(len(nameBuf.Bytes())))
		if err != nil {
			return fmt.Errorf("write object dataBuf %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.position.X)
		if err != nil {
			return fmt.Errorf("write object pos x %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.position.Y)
		if err != nil {
			return fmt.Errorf("write object pos y %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.position.Z)
		if err != nil {
			return fmt.Errorf("write object pos z %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.rotation.X)
		if err != nil {
			return fmt.Errorf("write object rot x %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.rotation.Y)
		if err != nil {
			return fmt.Errorf("write object rot y %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.rotation.Z)
		if err != nil {
			return fmt.Errorf("write object rot z %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.scale)
		if err != nil {
			return fmt.Errorf("write object scale %s: %w", o.name, err)
		}
	}

	for _, o := range e.regions {
		names = append(names, o.name)
		err = helper.WriteString(nameBuf, o.name)
		if err != nil {
			return fmt.Errorf("write region nameBuf %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(len(nameBuf.Bytes())))
		if err != nil {
			return fmt.Errorf("write region dataBuf %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.center.X)
		if err != nil {
			return fmt.Errorf("write region center x %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.center.Y)
		if err != nil {
			return fmt.Errorf("write region center y %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.center.Z)
		if err != nil {
			return fmt.Errorf("write region center z %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.unknown.X)
		if err != nil {
			return fmt.Errorf("write region unknown a %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.unknown.Y)
		if err != nil {
			return fmt.Errorf("write region unknown b %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.unknown.Z)
		if err != nil {
			return fmt.Errorf("write region unknown c %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.extent.X)
		if err != nil {
			return fmt.Errorf("write region extent x %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.extent.Y)
		if err != nil {
			return fmt.Errorf("write region extent y %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.extent.Z)
		if err != nil {
			return fmt.Errorf("write region extent z %s: %w", o.name, err)
		}
	}

	for _, o := range e.lights {
		names = append(names, o.name)
		err = helper.WriteString(nameBuf, o.name)
		if err != nil {
			return fmt.Errorf("write light nameBuf %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, uint32(len(nameBuf.Bytes())))
		if err != nil {
			return fmt.Errorf("write light dataBuf %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.position.X)
		if err != nil {
			return fmt.Errorf("write light position x %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.position.Y)
		if err != nil {
			return fmt.Errorf("write light position y %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.position.Z)
		if err != nil {
			return fmt.Errorf("write light position z %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.color.R)
		if err != nil {
			return fmt.Errorf("write light color r %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.color.G)
		if err != nil {
			return fmt.Errorf("write light color g %s: %w", o.name, err)
		}
		err = binary.Write(dataBuf, binary.LittleEndian, o.color.B)
		if err != nil {
			return fmt.Errorf("write light color b %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.radius)
		if err != nil {
			return fmt.Errorf("write light radius %s: %w", o.name, err)
		}
	}

	// Start writing
	err = binary.Write(w, binary.LittleEndian, []byte("EQGZ"))
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(1))
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(names)))
	if err != nil {
		return fmt.Errorf("write name count: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(e.models)))
	if err != nil {
		return fmt.Errorf("write model count: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(e.objects)))
	if err != nil {
		return fmt.Errorf("write object count: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(e.regions)))
	if err != nil {
		return fmt.Errorf("write region count: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(e.lights)))
	if err != nil {
		return fmt.Errorf("write light count: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, nameBuf.Bytes())
	if err != nil {
		return fmt.Errorf("write nameBuf: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, dataBuf.Bytes())
	if err != nil {
		return fmt.Errorf("write dataBuf: %w", err)
	}
	return nil
}
