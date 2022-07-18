package gltf

import (
	"fmt"

	"github.com/qmuntal/gltf/modeler"

	"github.com/qmuntal/gltf"
)

type Primitive struct {
	MaterialIndex *uint32
	Positions     [][3]float32
	Normals       [][3]float32
	Uvs           [][2]float32
	Joints        [][4]uint16
	Weights       [][4]float32
	Indices       []uint16
	UniqueIndices map[uint32]uint16
}

func NewPrimitive() *Primitive {
	return &Primitive{
		UniqueIndices: make(map[uint32]uint16),
	}
}

func (e *GLTF) PrimitiveAdd(meshName string, prim *Primitive) error {
	primitive := &gltf.Primitive{
		Mode:     gltf.PrimitiveTriangles,
		Material: prim.MaterialIndex,
	}
	primitive.Attributes = map[string]uint32{
		gltf.POSITION:   modeler.WritePosition(e.doc, prim.Positions),
		gltf.NORMAL:     modeler.WriteNormal(e.doc, prim.Normals),
		gltf.TEXCOORD_0: modeler.WriteTextureCoord(e.doc, prim.Uvs),
	}
	if len(prim.Joints) > 0 {
		primitive.Attributes[gltf.JOINTS_0] = modeler.WriteJoints(e.doc, prim.Joints)
		primitive.Attributes[gltf.WEIGHTS_0] = modeler.WriteWeights(e.doc, prim.Weights)
	}

	primitive.Indices = gltf.Index(modeler.WriteIndices(e.doc, prim.Indices))

	entry := e.meshes[meshName]
	if entry == nil {
		return fmt.Errorf("mesh %s not found", meshName)
	}

	entry.mesh.Primitives = append(entry.mesh.Primitives, primitive)
	return nil
}
