package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/xackery/encdec"
)

// TriangleBuild prepares an EQG-styled triangle buffer list
func TriangleBuild(version uint32, names map[string]int32, matManager *MaterialManager, meshManager *MeshManager) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	if len(meshManager.Meshes()) < 1 {
		return nil, fmt.Errorf("no meshes found")
	}
	mesh := meshManager.meshes[0]
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)

	// triangles
	for _, o := range mesh.Triangles {
		nameOffset := int32(-1)
		for key, val := range names {
			if key == o.MaterialName {
				nameOffset = val
				break
			}
		}
		if nameOffset == -1 {
			//log.Debugf("material %s not found ignoring", o.MaterialName)
		}
		enc.Uint32(o.Index.X)
		enc.Uint32(o.Index.Y)
		enc.Uint32(o.Index.Z)
		enc.Int32(nameOffset)
		enc.Uint32(o.Flag)
	}
	return dataBuf.Bytes(), nil
}
