package geo

import "fmt"

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

// Vector3 has X,Y,Z defined as float32
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func NewVector3() *Vector3 {
	return &Vector3{}
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

func NewQuad4() *Quad4 {
	return &Quad4{}
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

func NewIndex3() *Index3 {
	return &Index3{}
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

func NewUIndex3() *UIndex3 {
	return &UIndex3{}
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

func NewIndex4() *Index4 {
	return &Index4{}
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

func NewTriangle() *Triangle {
	return &Triangle{
		Index: &UIndex3{},
	}
}

// String returns a string version of triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("%s %s %d", t.Index, t.MaterialName, t.Flag)
}
