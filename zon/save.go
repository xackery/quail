package zon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Save writes a zon file to location
func (e *ZON) Save(w io.Writer) error {
	var err error

	type nameInfo struct {
		offset uint32
		name   string
	}
	names := []*nameInfo{}

	nameBuf := bytes.NewBuffer(nil)
	dataBuf := bytes.NewBuffer(nil)

	for _, o := range e.models {
		name := &nameInfo{
			name:   o.name,
			offset: uint32(nameBuf.Len()),
		}

		isNew := true
		for _, n := range names {
			if n.name != name.name {
				continue
			}
			isNew = false
			name = n
			break
		}
		if isNew {
			names = append(names, name)
			err = binary.Write(nameBuf, binary.LittleEndian, []byte(o.name))
			if err != nil {
				return fmt.Errorf("write name %s: %w", o.name, err)
			}
			err = binary.Write(nameBuf, binary.LittleEndian, []byte{0})
			if err != nil {
				return fmt.Errorf("write zero %s: %w", o.name, err)
			}
		}

		err = binary.Write(dataBuf, binary.LittleEndian, name.offset)
		if err != nil {
			return fmt.Errorf("write name offset %s: %w", o.name, err)
		}
	}

	for _, o := range e.objects {

		modelID := uint32(9999)
		for i := range names {
			if names[i].name != o.modelName {
				continue
			}
			modelID = uint32(i)
			break
		}
		if modelID == 9999 {
			return fmt.Errorf("modelID %s not found", o.modelName)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, modelID)
		if err != nil {
			return fmt.Errorf("write object model id %s: %w", o.name, err)
		}
		//binary.Write(dataBuf, binary.LittleEndian, uint16(0))

		name := &nameInfo{
			name:   o.name,
			offset: uint32(nameBuf.Len()),
		}

		isNew := true
		for _, n := range names {
			if n.name != name.name {
				continue
			}
			isNew = false
			name = n
			break
		}
		if isNew {
			names = append(names, name)
			err = binary.Write(nameBuf, binary.LittleEndian, []byte(o.name))
			if err != nil {
				return fmt.Errorf("write name %s: %w", o.name, err)
			}
			err = binary.Write(nameBuf, binary.LittleEndian, []byte{0})
			if err != nil {
				return fmt.Errorf("write zero %s: %w", o.name, err)
			}
		}

		err = binary.Write(dataBuf, binary.LittleEndian, name.offset)
		if err != nil {
			return fmt.Errorf("write objectname offset %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.position)
		if err != nil {
			return fmt.Errorf("write object pos %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.rotation)
		if err != nil {
			return fmt.Errorf("write object rot %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.scale)
		if err != nil {
			return fmt.Errorf("write object scale %s: %w", o.name, err)
		}

	}

	for _, o := range e.regions {
		name := &nameInfo{
			name:   o.name,
			offset: uint32(nameBuf.Len()),
		}

		isNew := true
		for _, n := range names {
			if n.name != name.name {
				continue
			}
			isNew = false
			name = n
			break
		}
		if isNew {
			names = append(names, name)
			err = binary.Write(nameBuf, binary.LittleEndian, []byte(o.name))
			if err != nil {
				return fmt.Errorf("write name %s: %w", o.name, err)
			}
			err = binary.Write(nameBuf, binary.LittleEndian, []byte{0})
			if err != nil {
				return fmt.Errorf("write zero %s: %w", o.name, err)
			}
		}

		err = binary.Write(dataBuf, binary.LittleEndian, name.offset)
		if err != nil {
			return fmt.Errorf("write region name offset %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.center)
		if err != nil {
			return fmt.Errorf("write region center %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.unknown)
		if err != nil {
			return fmt.Errorf("write region unknown %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.extent)
		if err != nil {
			return fmt.Errorf("write region extent %s: %w", o.name, err)
		}
	}

	for _, o := range e.lights {
		name := &nameInfo{
			name:   o.name,
			offset: uint32(nameBuf.Len()),
		}

		isNew := true
		for _, n := range names {
			if n.name != name.name {
				continue
			}
			isNew = false
			name = n
			break
		}
		if isNew {
			names = append(names, name)
			err = binary.Write(nameBuf, binary.LittleEndian, []byte(o.name))
			if err != nil {
				return fmt.Errorf("write name %s: %w", o.name, err)
			}
			err = binary.Write(nameBuf, binary.LittleEndian, []byte{0})
			if err != nil {
				return fmt.Errorf("write zero %s: %w", o.name, err)
			}
		}

		err = binary.Write(dataBuf, binary.LittleEndian, name.offset)
		if err != nil {
			return fmt.Errorf("write light name offset %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.position)
		if err != nil {
			return fmt.Errorf("write light position %s: %w", o.name, err)
		}

		err = binary.Write(dataBuf, binary.LittleEndian, o.color)
		if err != nil {
			return fmt.Errorf("write light color %s: %w", o.name, err)
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
	err = binary.Write(w, binary.LittleEndian, uint32(nameBuf.Len()))
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
