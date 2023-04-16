package geo

import (
	"fmt"
	"io"
)

// Bone is a bone
type Bone struct {
	Name          string
	Next          int32
	ChildrenCount uint32
	ChildIndex    int32
	Pivot         Vector3
	Rotation      Quad4
	Scale         Vector3
}

func NewBone() Bone {
	return Bone{
		Pivot:    NewVector3(),
		Rotation: NewQuad4(),
		Scale:    NewVector3(),
	}
}

// WriteHeader writes the header for a Bone
func (e *Bone) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("name|child_index|children_count|next|pivot|rotation|scale\n")
	return err
}

// WriteString writes a Bone to a string writer
func (e *Bone) Write(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("%s|%d|%d|%d|%s|%s|%0.3f\n",
		e.Name,
		e.ChildIndex,
		e.ChildrenCount,
		e.Next,
		&Vector3{X: e.Pivot.Y, Y: -e.Pivot.X, Z: e.Pivot.Z},
		e.Rotation,
		e.Scale))
	return err
}
