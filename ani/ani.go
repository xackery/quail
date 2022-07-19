// ani are animation files, found in EverQuest eqg files
package ani

import "github.com/g3n/engine/math32"

type ANI struct {
	name     string
	bones    []*Bone
	isStrict bool
}

type Bone struct {
	frameCount  uint32
	name        string
	delay       int32
	translation *math32.Vector3
	rotation    *math32.Vector4
	scale       *math32.Vector3
}

func New(name string) (*ANI, error) {
	e := &ANI{
		name: name,
	}
	return e, nil
}
