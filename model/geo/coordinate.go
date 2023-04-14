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

// NewVector2 returns a new vector2
func NewVector2() *Vector2 {
	return &Vector2{}
}

// NewVector2FromString returns a new vector2 from a string
func NewVector2FromString(s string) *Vector2 {
	parts := strings.Split(s, "|")
	if len(parts) < 2 {
		return nil
	}
	return &Vector2{
		X: helper.AtoF32(parts[0]),
		Y: helper.AtoF32(parts[1]),
	}
}

// String returns a string version of vector2
func (v *Vector2) String() string {
	return fmt.Sprintf("%0.3f,%0.3f", v.X, v.Y)
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

// String returns a string version of vector3
func (e *Vector3) String() string {
	return fmt.Sprintf("%0.3f,%0.3f,%0.3f", e.X, e.Y, e.Z)
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
	return fmt.Sprintf("%0.3f,%0.3f,%0.3f,%0.3f", q.X, q.Y, q.Z, q.W)
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
	return fmt.Sprintf("%d|%d|%d|%d", i.X, i.Y, i.Z, i.W)
}
