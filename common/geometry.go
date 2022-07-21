package common

import (
	"fmt"

	"github.com/g3n/engine/math32"
)

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
	Tint     *Tint
	Uv       *math32.Vector2
	Uv2      *math32.Vector2
}

type Triangle struct {
	Index        [3]uint32
	MaterialName string
	Flag         uint32
}

type Tint struct {
	R uint8
	G uint8
	B uint8
}
