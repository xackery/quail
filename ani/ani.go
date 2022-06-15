package ani

type ANI struct {
	name  string
	bones []*Bone
}

type Bone struct {
	Delay       int32
	Translation [3]float32
	Rotation    [4]float32
	Scale       [3]float32
}

func New(name string) (*ANI, error) {
	e := &ANI{
		name: name,
	}
	return e, nil
}
