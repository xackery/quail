package mod

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// Encode writes a zon file to location
func (e *MOD) Encode(w io.Writer) error {
	var err error

	nameData, data, err := geo.WriteGeometry(e.version, e.MaterialManager, e.meshManager)
	if err != nil {
		return fmt.Errorf("writeGeometry: %w", err)
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGM")
	enc.Uint32(e.version)
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(e.MaterialManager.Count()))
	enc.Uint32(uint32(e.meshManager.VertexTotalCount()))
	enc.Uint32(uint32(e.meshManager.TriangleTotalCount()))
	enc.Uint32(uint32(e.meshManager.BoneTotalCount()))
	enc.Bytes(nameData)
	enc.Bytes(data)
	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	log.Debugf("%s encoded %d verts, %d triangles, %d bones, %d materials", e.name, e.meshManager.VertexTotalCount(), e.meshManager.TriangleTotalCount(), e.meshManager.BoneTotalCount(), e.MaterialManager.Count())
	return nil
}
