package gltf

import (
	"fmt"

	"github.com/qmuntal/gltf/modeler"

	"github.com/qmuntal/gltf"
)

type primitiveEntry struct {
	index     *uint32
	primitive *gltf.Primitive
}

type Primitive struct {
	MaterialIndex *uint32
	Positions     [][3]float32
	Normals       [][3]float32
	Uvs           [][2]float32
	//joints        [][4]uint16
	//weights       [][4]uint16
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
		//gltf.JOINTS_0:   modeler.WriteJoints(doc, prim.Joints),
		//gltf.WEIGHTS_0:  modeler.WriteWeights(doc, prim.Weights),
	}

	primitive.Indices = gltf.Index(modeler.WriteIndices(e.doc, prim.Indices))

	entry := e.meshes[meshName]
	if entry == nil {
		return fmt.Errorf("mesh %s not found", meshName)
	}

	entry.mesh.Primitives = append(entry.mesh.Primitives, primitive)
	return nil
}
