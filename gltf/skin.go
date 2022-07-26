package gltf

import (
	"github.com/qmuntal/gltf"
)

func (e *GLTF) SkinAdd(skin *gltf.Skin) uint32 {
	e.doc.Skins = append(e.doc.Skins, skin)
	return uint32(len(e.doc.Skins) - 1)
}
