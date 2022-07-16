package mds

import (
	"github.com/g3n/engine/math32"
)

func (e *MDS) BoneAdd(name string, next int32, childrenCount uint32, childIndex int32, pivot *math32.Vector3, rot *math32.Vector4, scale *math32.Vector3) error {
	bone := &bone{
		name:          name,
		next:          next,
		childrenCount: childrenCount,
		childIndex:    childIndex,
		pivot:         pivot,
		rot:           rot,
		scale:         scale,
	}

	e.bones = append(e.bones, bone)
	return nil
}
