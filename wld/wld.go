package wld

import (
	"fmt"

	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wld/animation"
	"github.com/xackery/quail/wld/bsp"
	"github.com/xackery/quail/wld/material"
	"github.com/xackery/quail/wld/mesh"
)

type Wld struct {
	Identifier   int
	Version      int
	bspTree      []*bsp.BspTree
	fragments    map[int]fragBase
	skeletons    []*animation.SkeletonHierarchy
	trackDefs    []*animation.TrackDefFragment
	tracks       []*animation.TrackFragment
	meshes       []*mesh.Mesh
	materialList []*material.MaterialList
	actors       []*animation.ActorInstance
	objects      []*animation.ActorDef
}

type fragBase interface {
	FragCode() int
}

func (e *Wld) processFragments(fragments map[int]raw.FragmentReadWriter) error {
	e.fragments = make(map[int]fragBase)

	if len(fragments) == 0 {
		return fmt.Errorf("no fragments found")
	}
	maxFragments := len(fragments)
	for i := 1; i < maxFragments; i++ {
		frag, ok := fragments[i]
		if !ok {
			return fmt.Errorf("fragment %d not found", i)
		}
		switch frag.FragCode() {
		case 0x03: // Texture Path
			//e.fragments[i] = &raw.
		}
	}
	return nil
}
