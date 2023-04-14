package geo

import (
	"fmt"
	"io"
)

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

// WriteHeader writes the header for a BoneAnimation
func (e *BoneAnimation) WriteHeader(w io.StringWriter) error {
	_, err := w.WriteString("name|frame_count\n")
	return err
}

// Write writes a BoneAnimation to a string writer
func (e *BoneAnimation) Write(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("%s|%d\n", e.Name, e.FrameCount))
	return err
}
