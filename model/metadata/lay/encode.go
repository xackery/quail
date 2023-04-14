package lay

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type nameInfo struct {
	offset uint32
	name   string
}

// Encode writes a zon file to location
func (e *LAY) Encode(w io.Writer) error {
	var err error

	//prepare name buffer
	names := []*nameInfo{}
	nameBuf := bytes.NewBuffer(nil)
	// materials

	tmpNames := []string{}
	for _, o := range e.layerManager.Layers() {
		tmpNames = append(tmpNames, o.Entry0)
		tmpNames = append(tmpNames, o.Entry1)
	}

	for _, name := range tmpNames {
		isNew := true
		for _, val := range names {
			if val.name == name {
				isNew = false
				break
			}
		}
		if !isNew {
			continue
		}

		names = append(names, &nameInfo{
			offset: uint32(nameBuf.Len()),
			name:   name,
		})
		_, err = nameBuf.Write([]byte(name))
		if err != nil {
			return fmt.Errorf("write name: %w", err)
		}
		_, err = nameBuf.Write([]byte{0})
		if err != nil {
			return fmt.Errorf("write 0: %w", err)
		}
	}

	nameData := nameBuf.Bytes()

	// Start writing
	err = binary.Write(w, binary.LittleEndian, []byte("EQGL"))
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(1))
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}
	//versionOffset := 32
	err = binary.Write(w, binary.LittleEndian, uint32(len(nameData)))
	if err != nil {
		return fmt.Errorf("write name length: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(e.layerManager.Count()))
	if err != nil {
		return fmt.Errorf("write materialCount: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, nameData)
	if err != nil {
		return fmt.Errorf("write nameBuf: %w", err)
	}
	/*
		for _, o := range e.layers {
			err = binary.Write(w, binary.LittleEndian, o.name)
			if err != nil {
				return fmt.Errorf("write dataBuf: %w", err)
			}
		}*/
	return fmt.Errorf("not supported")
	//return nil
}
