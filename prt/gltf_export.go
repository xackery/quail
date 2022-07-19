package prt

import (
	"fmt"

	qgltf "github.com/xackery/quail/gltf"
)

// GLTFExport exports a provided prt file to gltf format
func (e *PRT) GLTFExport(doc *qgltf.GLTF) error {
	var err error
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}
	for _, entry := range e.particles {
		err = doc.PaticleAdd(entry)
		if err != nil {
			return fmt.Errorf("particleAdd: %w", err)
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
