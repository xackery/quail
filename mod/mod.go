package mod

import "github.com/g3n/engine/math32"

// MOD is a zon file struct
type MOD struct {
	materials       []*material
	vertices        []*vertex
	triangles       []*triangle
	bones           []*bone
	boneAssignments []*boneAssignment
}

type material struct {
	name       string
	shaderName string
	properties []*property
}

type property struct {
	name       string
	typeValue  uint32
	floatValue float32
	intValue   uint32
}

type vertex struct {
	position  math32.Vector3
	position2 math32.Vector3
	uv        math32.Vector2
}

type triangle struct {
	index        math32.Vector3
	materialName string
	flag         uint32
}

type bone struct {
	name    string
	unknown [13]float32
}

type boneAssignment struct {
	unknown [9]uint32
}
