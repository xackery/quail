package geo

// Bone is a bone
type Bone struct {
	Name          string
	Next          int32
	ChildrenCount uint32
	ChildIndex    int32
	Pivot         *Vector3
	Rotation      *Quad4
	Scale         *Vector3
}

func NewBone() *Bone {
	return &Bone{
		Pivot:    NewVector3(),
		Rotation: NewQuad4(),
		Scale:    NewVector3(),
	}
}
