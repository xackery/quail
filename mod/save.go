package mod

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Save writes a zon file to location
func (e *MOD) Save(w io.Writer) error {
	var err error

	nameData, data, err := e.writeGeometry()
	if err != nil {
		return fmt.Errorf("writeGeometry: %w", err)
	}

	// Start writing
	err = binary.Write(w, binary.LittleEndian, []byte("EQGM"))
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(1))
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(nameData)))
	if err != nil {
		return fmt.Errorf("write name length: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(e.materials)))
	if err != nil {
		return fmt.Errorf("write material count: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(e.vertices)))
	if err != nil {
		return fmt.Errorf("write vertices count: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(e.triangles)))
	if err != nil {
		return fmt.Errorf("write triangle count: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(e.bones)))
	if err != nil {
		return fmt.Errorf("write bone count: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(len(e.boneAssignments)))
	if err != nil {
		return fmt.Errorf("write bone assignemt count: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, nameData)
	if err != nil {
		return fmt.Errorf("write nameBuf: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, data)
	if err != nil {
		return fmt.Errorf("write dataBuf: %w", err)
	}
	return nil
}
