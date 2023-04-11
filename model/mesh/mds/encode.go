package mds

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/model/geo"
)

// Encode writes a zon file to location
func (e *MDS) Encode(w io.Writer) error {
	var err error

	nameData, data, err := geo.WriteGeometry(e.version, e.materials, e.vertices, e.triangles, e.bones)
	if err != nil {
		return fmt.Errorf("writeGeometry: %w", err)
	}

	// Start writing
	err = binary.Write(w, binary.LittleEndian, []byte("EQGS"))
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

	err = binary.Write(w, binary.LittleEndian, uint32(len(e.bones)))
	if err != nil {
		return fmt.Errorf("write bone count: %w", err)
	}

	//TODO: mds encode sub count
	err = binary.Write(w, binary.LittleEndian, uint32(0))
	if err != nil {
		return fmt.Errorf("write sub count: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, nameData)
	if err != nil {
		return fmt.Errorf("write nameBuf: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, data)
	if err != nil {
		return fmt.Errorf("write dataBuf: %w", err)
	}

	//TODO: bone data

	return nil
}
