package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/xackery/encdec"
)

// VertexBuild prepares an EQG-styled vertex buffer list
func VertexBuild(version uint32, names map[string]int32, meshManager *MeshManager) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	if len(meshManager.Meshes()) < 1 {
		return nil, fmt.Errorf("no meshes found")
	}
	mesh := meshManager.meshes[0]
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)
	// verts
	for _, o := range mesh.Vertices {
		enc.Float32(o.Position.X)
		enc.Float32(o.Position.Y)
		enc.Float32(o.Position.Z)
		enc.Float32(o.Normal.X)
		enc.Float32(o.Normal.Y)
		enc.Float32(o.Normal.Z)
		if version < 3 {
			enc.Float32(o.Uv.X)
			enc.Float32(o.Uv.Y)
		} else {
			enc.Uint8(o.Tint.R)
			enc.Uint8(o.Tint.G)
			enc.Uint8(o.Tint.B)
			enc.Uint8(o.Tint.A)
			enc.Float32(o.Uv.X)
			enc.Float32(o.Uv.Y)
			enc.Float32(o.Uv2.X)
			enc.Float32(o.Uv2.Y)
		}
	}
	return dataBuf.Bytes(), nil
}
