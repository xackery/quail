package geo

import "fmt"

// BoneAnimation is a bone animation
type BoneAnimation struct {
	Name       string
	FrameCount uint32
	Frames     []*BoneAnimationFrame
}

// String returns a string representation of the bone animation
func (e *BoneAnimation) String() string {
	return fmt.Sprintf("%s %d", e.Name, e.FrameCount)
}

// BoneAnimationFrame is a bone animation frame
type BoneAnimationFrame struct {
	Milliseconds uint32
	Translation  *Vector3
	Rotation     *Quad4
	Scale        *Vector3
}

// String returns a string representation of the bone frame
func (e *BoneAnimationFrame) String() string {
	return fmt.Sprintf("%d %s %s %s", e.Milliseconds, e.Translation, e.Rotation, e.Scale)
}
