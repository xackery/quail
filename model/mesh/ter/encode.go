package ter

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// Encode writes a zon file to location
func (e *TER) Encode(w io.Writer) error {
	var err error
	if e.version != 2 {
		log.Warnf("%s: version %d not supported, using 2", e.name, e.version)
		e.version = 2
	}

	nameData, meshData, err := geo.WriteGeometry(e.version, e.MaterialManager, e.meshManager)
	if err != nil {
		return fmt.Errorf("writeGeometry: %w", err)
	}

	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGT")
	enc.Uint32(e.version)
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(e.MaterialManager.Count()))
	enc.Uint32(uint32(e.meshManager.VertexCount(e.name)))
	enc.Uint32(uint32(e.meshManager.TriangleCount(e.name)))
	enc.Bytes(nameData)
	enc.Bytes(meshData)

	log.Debugf("%s encoded %d verts, %d triangles, %d materials", e.name, e.meshManager.VertexTotalCount(), e.meshManager.TriangleTotalCount(), e.MaterialManager.Count())
	return nil
}
