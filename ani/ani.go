// ani are animation files, found in EverQuest eqg files
package ani

type ANI struct {
	name     string
	bones    []*Bone
	isStrict bool
}

type Bone struct {
	frameCount  uint32
	name        string
	delay       int32
	translation [3]float32
	rotation    [4]float32
	scale       [3]float32
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string) (*ANI, error) {
	e := &ANI{
		name: name,
	}
	return e, nil
}
