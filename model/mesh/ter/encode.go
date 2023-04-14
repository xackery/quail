package ter

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/model/geo"
)

// Encode writes a zon file to location
func (e *TER) Encode(w io.Writer) error {
	var err error

	nameData, data, err := geo.WriteGeometry(e.version, e.MaterialManager, e.meshManager)
	if err != nil {
		return fmt.Errorf("writeGeometry: %w", err)
	}

	// Start writing
	err = binary.Write(w, binary.LittleEndian, []byte("EQGT"))
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(2))
	if err != nil {
		return fmt.Errorf("write header version: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(len(nameData)))
	if err != nil {
		return fmt.Errorf("write name length: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(e.MaterialManager.Count()))
	if err != nil {
		return fmt.Errorf("write material count: %w", err)
	}
	err = binary.Write(w, binary.LittleEndian, uint32(e.meshManager.VertexCount(e.name)))
	if err != nil {
		return fmt.Errorf("write vertices count: %w", err)
	}

	err = binary.Write(w, binary.LittleEndian, uint32(e.meshManager.TriangleCount(e.name)))
	if err != nil {
		return fmt.Errorf("write triangle count: %w", err)
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
