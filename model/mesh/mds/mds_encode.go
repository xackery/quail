package mds

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

// Encode writes a mds file
func Encode(mesh *common.Model, version uint32, w io.Writer) error {
	var err error
	modelNames := []string{}

	if len(mesh.Bones) == 0 {
		mesh.Bones = append(mesh.Bones, common.Bone{
			Name:       "ROOT_BONE",
			Next:       -1,
			ChildIndex: -1,
			Scale:      common.Vector3{X: 1, Y: 1, Z: 1},
		})
	}

	if len(mesh.Bones) > 0 {
		modelNames = append(modelNames, mesh.Name)
	}

	names, nameData, err := mesh.NameBuild(modelNames)
	if err != nil {
		return fmt.Errorf("nameBuild: %w", err)
	}

	materialData, err := mesh.MaterialBuild(names)
	if err != nil {
		return fmt.Errorf("materialBuild: %w", err)
	}

	verticesData, err := mesh.VertexBuild(version, names)
	if err != nil {
		return fmt.Errorf("vertexBuild: %w", err)
	}

	triangleData, err := mesh.TriangleBuild(version, names)
	if err != nil {
		return fmt.Errorf("triangleBuild: %w", err)
	}

	boneData, err := mesh.BoneBuild(version, "mds", names)
	if err != nil {
		return fmt.Errorf("boneBuild: %w", err)
	}

	tag.New()
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGS")
	enc.Uint32(version)
	enc.Uint32(uint32(len(nameData)))
	enc.Uint32(uint32(len(mesh.Materials)))
	enc.Uint32(uint32(len(mesh.Bones)))
	enc.Uint32(0) // subCount
	enc.Bytes(nameData)
	enc.Bytes(materialData)
	enc.Bytes(boneData)
	enc.Uint32(0) // mainNameIndex?
	enc.Uint32(0) // subNameIndex?
	enc.Uint32(uint32(len(mesh.Vertices)))
	enc.Uint32(uint32(len(mesh.Triangles)))
	enc.Uint32(0) //TODO: fix boneassignmentcount
	enc.Bytes(verticesData)
	enc.Bytes(triangleData)

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	log.Debugf("%s encoded %d verts, %d triangles, %d bones, %d materials", mesh.Name, len(mesh.Vertices), len(mesh.Triangles), len(mesh.Bones), len(mesh.Materials))
	return nil
}
