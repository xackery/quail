package ter

import "github.com/g3n/engine/math32"

// TER is a zon file struct
type TER struct {
	materials []*material
	vertices  []*vertex
	triangles []*triangle
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
