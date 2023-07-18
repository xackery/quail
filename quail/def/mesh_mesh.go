package def

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/xackery/encdec"
)

// VertexBuild prepares an EQG-styled vertex buffer list
func (mesh *Mesh) vertexBuild(version uint32, names map[string]int32) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
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

// TriangleBuild prepares an EQG-styled triangle buffer list
func (mesh *Mesh) triangleBuild(version uint32, names map[string]int32) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)

	// triangles
	for _, o := range mesh.Triangles {
		materialIdx := int32(-1)
		for idx, val := range mesh.Materials {
			if val.Name != o.MaterialName {
				continue
			}
			materialIdx = int32(idx)
			break
		}
		enc.Uint32(o.Index.X)
		enc.Uint32(o.Index.Y)
		enc.Uint32(o.Index.Z)
		enc.Int32(materialIdx)
		enc.Uint32(o.Flag)
	}
	return dataBuf.Bytes(), nil
}

// BoneBuild prepares an EQG-styled bone buffer list
func (mesh *Mesh) boneBuild(version uint32, isMod bool, names map[string]int32) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)

	// bones
	for _, o := range mesh.Bones {
		nameOffset := int32(-1)
		for key, val := range names {
			if key == o.Name {
				nameOffset = val
				break
			}
		}
		if nameOffset == -1 {
			return nil, fmt.Errorf("bone %s not found", o.Name)
		}

		enc.Int32(nameOffset)
		enc.Int32(o.Next)
		enc.Uint32(o.ChildrenCount)
		enc.Int32(o.ChildIndex)
		enc.Float32(o.Pivot.X)
		enc.Float32(o.Pivot.Y)
		enc.Float32(o.Pivot.Z)
		enc.Float32(o.Rotation.X)
		enc.Float32(o.Rotation.Y)
		enc.Float32(o.Rotation.Z)
		//enc.Float32(o.Rotation.W)
		enc.Float32(o.Scale.X)
		enc.Float32(o.Scale.Y)
		enc.Float32(o.Scale.Z)
		if isMod {
			enc.Float32(1.0)
		}
	}
	return dataBuf.Bytes(), nil
}
