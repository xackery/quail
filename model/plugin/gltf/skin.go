package gltf

import (
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func (e *GLTF) SkinAdd(skin *gltf.Skin) uint32 {
	e.doc.Skins = append(e.doc.Skins, skin)
	return uint32(len(e.doc.Skins) - 1)
}

func (e *GLTF) WriteMatrix(matrix [][4][4]float32) uint32 {
	return modeler.WriteAccessor(e.doc, gltf.TargetArrayBuffer, matrix)
}
