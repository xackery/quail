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
	names, nameData, err := geo.NameBuild(e.MaterialManager, e.meshManager)
	if err != nil {
		return fmt.Errorf("nameBuild: %w", err)
	}

	materialData, err := geo.MaterialBuild(e.version, names, e.MaterialManager)
	if err != nil {
		return fmt.Errorf("materialBuild: %w", err)
	}

	verticesData, err := geo.VertexBuild(e.version, names, e.meshManager)
	if err != nil {
		return fmt.Errorf("vertexBuild: %w", err)
	}

	triangleData, err := geo.TriangleBuild(e.version, names, e.MaterialManager, e.meshManager)
	if err != nil {
		return fmt.Errorf("triangleBuild: %w", err)
	}

	boneData, err := geo.BoneBuild(e.version, names, e.meshManager)
	if err != nil {
		return fmt.Errorf("boneBuild: %w", err)
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
	enc.Bytes(materialData)
	enc.Bytes(verticesData)
	enc.Bytes(triangleData)
	enc.Bytes(boneData)

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	log.Debugf("%s encoded %d verts, %d triangles, %d bones, %d materials", e.name, e.meshManager.VertexTotalCount(), e.meshManager.TriangleTotalCount(), e.meshManager.BoneTotalCount(), e.MaterialManager.Count())
	return nil
}
