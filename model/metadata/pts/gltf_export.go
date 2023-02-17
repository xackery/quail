package pts

import (
	"fmt"

	qgltf "github.com/xackery/quail/model/plugin/gltf"
)

// GLTFEncode exports a provided pts file to gltf format
func (e *PTS) GLTFEncode(doc *qgltf.GLTF) error {
	var err error
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}
	for _, entry := range e.particles {
		err = doc.ParticlePointAdd(entry)
		if err != nil {
			return fmt.Errorf("ParticlePointAdd: %w", err)
		}
	}
	return nil
}

/*
func (e *MOD) gltfBoneChildren(doc *gltf.Document, children *[]uint32, boneIndex int) error {

	nodeIndex, ok := e.gltfBoneBuffer[boneIndex]
	if !ok {
		return fmt.Errorf("bone %d node not found", boneIndex)
	}
	*children = append(*children, nodeIndex)

	bone := e.bones[boneIndex]
	if bone.next == -1 {
		return nil
	}

	return e.gltfBoneChildren(doc, children, int(bone.next))
}*/
