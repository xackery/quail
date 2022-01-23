package common

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

type Material struct {
	Name       string
	ShaderName string
	Flag       uint32
	Properties []*Property
}

func (e *Material) String() string {
	return fmt.Sprintf("{Name: %s, Flag: %d, Properties: (%d)}", e.Name, e.Flag, len(e.Properties))
}

type Property struct {
	Name       string
	TypeValue  uint32
	FloatValue float32
	IntValue   uint32
	StrValue   string
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
