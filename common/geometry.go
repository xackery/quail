package common

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

type Material struct {
	Name       string
	ShaderName string
	Flag       uint32
	Properties Properties
}

func (e *Material) String() string {
	return fmt.Sprintf("{Name: %s, Flag: %d, Properties: (%s)}", e.Name, e.Flag, e.Properties)
}

type Properties []*Property

func (e Properties) String() string {
	value := ""
	for _, p := range e {
		value += fmt.Sprintf("[%s: %s], ", p.Name, p.Value)
	}
	if len(e) > 0 {
		value = value[0 : len(value)-2]
	}
	return value
}

type Property struct {
	Name     string
	Category uint32
	Value    string
}

type Vertex struct {
	Position *math32.Vector3
	Normal   *math32.Vector3
	Uv       *math32.Vector2
}

type Face struct {
	Index        [3]uint32
	MaterialName string
	Flag         uint32
}
