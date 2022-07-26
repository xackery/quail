package common

type Bone struct {
	Name          string
	Next          int32
	ChildrenCount uint32
	ChildIndex    int32
	Pivot         [3]float32
	Rotation      [4]float32
	Scale         [3]float32
}
