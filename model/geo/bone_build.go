package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/xackery/encdec"
)

// BoneBuild prepares an EQG-styled bone buffer list
func BoneBuild(version uint32, isMod bool, names map[string]int32, meshManager *MeshManager) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	if len(meshManager.Meshes()) < 1 {
		return nil, fmt.Errorf("no meshes found")
	}
	mesh := meshManager.meshes[0]
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
		enc.Float32(o.Rotation.W)
		enc.Float32(o.Scale.X)
		enc.Float32(o.Scale.Y)
		enc.Float32(o.Scale.Z)
		if isMod {
			enc.Float32(1.0)
		}
	}
	return dataBuf.Bytes(), nil
}
