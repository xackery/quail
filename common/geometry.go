package common

import "github.com/g3n/engine/math32"

type Material struct {
	Name       string
	ShaderName string
	Properties []*Property
}

type Property struct {
	Name       string
	TypeValue  uint32
	FloatValue float32
	IntValue   uint32
}

type Vertex struct {
	Position math32.Vector3
	Normal   math32.Vector3
	Uv       math32.Vector2
}

type Triangle struct {
	Index        math32.Vector3
	MaterialName string
	Flag         uint32
}
