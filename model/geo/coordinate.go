package geo

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/helper"
)

// Vector2 has X,Y defined as float32
type Vector2 struct {
	X float32
	Y float32
}

func NewVector2() *Vector2 {
	return &Vector2{}
}

// String returns a string version of vector2
func (v *Vector2) String() string {
	return fmt.Sprintf("%0.2f,%0.2f", v.X, v.Y)
}

// AtoVector2 converts a string to a vector2
func AtoVector2(s string) *Vector2 {
	parts := strings.Split(s, ",")
	if len(parts) < 2 {
		return nil
	}
	return &Vector2{
		X: helper.AtoF32(parts[0]),
		Y: helper.AtoF32(parts[1]),
	}
}

// Vector3 has X,Y,Z defined as float32
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// NewVector3 returns a new vector3
func NewVector3() *Vector3 {
	return &Vector3{}
}

// AtoVector3 converts a string to a vector3
func AtoVector3(s string) *Vector3 {
	parts := strings.Split(s, ",")
	if len(parts) < 3 {
		return nil
	}
	return &Vector3{
		X: helper.AtoF32(parts[0]),
		Y: helper.AtoF32(parts[1]),
		Z: helper.AtoF32(parts[2]),
	}
}

// String returns a string version of vector3
func (v *Vector3) String() string {
	return fmt.Sprintf("%0.2f,%0.2f,%f", v.X, v.Y, v.Z)
}

// Quad4  has X,Y,Z,W defined as float32
type Quad4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

// NewQuad4 returns a new quad4
func NewQuad4() *Quad4 {
	return &Quad4{}
}

// AtoQuad4 converts a string to a quad4
func AtoQuad4(s string) *Quad4 {
	parts := strings.Split(s, ",")
	if len(parts) < 4 {
		return nil
	}
	return &Quad4{
		X: helper.AtoF32(parts[0]),
		Y: helper.AtoF32(parts[1]),
		Z: helper.AtoF32(parts[2]),
		W: helper.AtoF32(parts[3]),
	}
}

// String returns a string version of quad4
func (q *Quad4) String() string {
	return fmt.Sprintf("%0.2f,%0.2f,%0.2f,%0.2f", q.X, q.Y, q.Z, q.W)
}

// Index3 has X,Y,Z defined as int32
type Index3 struct {
	X int32
	Y int32
	Z int32
}

// NewIndex3 returns a new index3
func NewIndex3() *Index3 {
	return &Index3{}
}

// AtoIndex3 converts a string to a index3
func AtoIndex3(s string) *Index3 {
	parts := strings.Split(s, ",")
	if len(parts) < 3 {
		return nil
	}
	return &Index3{
		X: helper.AtoI32(parts[0]),
		Y: helper.AtoI32(parts[1]),
		Z: helper.AtoI32(parts[2]),
	}
}

// String returns a string version of index3
func (i *Index3) String() string {
	return fmt.Sprintf("%d,%d,%d", i.X, i.Y, i.Z)
}

// UIndex3 has X,Y,Z defined as uint32
type UIndex3 struct {
	X uint32
	Y uint32
	Z uint32
}

// NewUIndex3 returns a new uindex3
func NewUIndex3() *UIndex3 {
	return &UIndex3{}
}

// AtoUIndex3 converts a string to a uindex3
func AtoUIndex3(s string) *UIndex3 {
	parts := strings.Split(s, ",")
	if len(parts) < 3 {
		return nil
	}
	return &UIndex3{
		X: helper.AtoU32(parts[0]),
		Y: helper.AtoU32(parts[1]),
		Z: helper.AtoU32(parts[2]),
	}
}

// String returns a string version of uindex3
func (i *UIndex3) String() string {
	return fmt.Sprintf("%d,%d,%d", i.X, i.Y, i.Z)
}

// Index4 has X,Y,Z,W defined as int16
type Index4 struct {
	X int16
	Y int16
	Z int16
	W int16
}

// NewIndex4 returns a new index4
func NewIndex4() *Index4 {
	return &Index4{}
}

// AtoIndex4 converts a string to a index4
func AtoIndex4(s string) *Index4 {
	parts := strings.Split(s, ",")
	if len(parts) < 4 {
		return nil
	}
	return &Index4{
		X: helper.AtoI16(parts[0]),
		Y: helper.AtoI16(parts[1]),
		Z: helper.AtoI16(parts[2]),
		W: helper.AtoI16(parts[3]),
	}
}

// String returns a string version of index4
func (i *Index4) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", i.X, i.Y, i.Z, i.W)
}

// Property contains data about a material
type Property struct {
	Name     string
	Category uint32
	Value    string
}

// NewProperty returns a new property
func NewProperty() *Property {
	return &Property{}
}

// Vertex stores information related to a mesh
type Vertex struct {
	Position *Vector3
	Normal   *Vector3
	Tint     *RGBA
	Uv       *Vector2
	Uv2      *Vector2
}

// NewVertex returns a new vertex
func NewVertex() *Vertex {
	return &Vertex{
		Position: &Vector3{},
		Normal:   &Vector3{},
		Tint:     &RGBA{},
		Uv:       &Vector2{},
		Uv2:      &Vector2{},
	}
}

// String returns a string version of vertex
func (v *Vertex) String() string {
	return fmt.Sprintf("%s %s %s %s %s", v.Position, v.Normal, v.Tint, v.Uv, v.Uv2)
}

// Triangle refers to the index of 3 vertices and maps it to a flag and material
type Triangle struct {
	Index        *UIndex3
	MaterialName string
	Flag         uint32
}

// NewTriangle returns a new triangle
func NewTriangle() *Triangle {
	return &Triangle{
		Index: &UIndex3{},
	}
}

// String returns a string version of triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("%s %s %d", t.Index, t.MaterialName, t.Flag)
}
