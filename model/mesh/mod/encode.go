package mod

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/tag"
)

// Encode writes a zon file to location
func (e *MOD) Encode(w io.Writer) error {
	var err error
	modelNames := []string{}
	if e.meshManager.BoneTotalCount() > 0 {
		modelNames = append(modelNames, e.name)
	}
	names, nameData, err := geo.NameBuild(e.MaterialManager, e.meshManager, modelNames)
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

	boneData, err := geo.BoneBuild(e.version, true, names, e.meshManager)
	if err != nil {
		return fmt.Errorf("boneBuild: %w", err)
	}

	tag.New()
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGM")
	enc.Uint32(e.version)
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(e.MaterialManager.Count()))
	enc.Uint32(uint32(e.meshManager.VertexTotalCount()))
	enc.Uint32(uint32(e.meshManager.TriangleTotalCount()))
	enc.Uint32(uint32(e.meshManager.BoneTotalCount()))
	tag.Add(0, int(enc.Pos()-1), "red", "header")
	enc.Bytes(nameData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "green", "names")
	enc.Bytes(materialData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "blue", "materials")
	enc.Bytes(verticesData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "yellow", "vertices")
	enc.Bytes(triangleData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "purple", "triangles")
	enc.Bytes(boneData)
	tag.Add(tag.LastPos(), int(enc.Pos()), "orange", "bones")

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	log.Debugf("%s encoded %d verts, %d triangles, %d bones, %d materials", e.name, e.meshManager.VertexTotalCount(), e.meshManager.TriangleTotalCount(), e.meshManager.BoneTotalCount(), e.MaterialManager.Count())
	return nil
}
