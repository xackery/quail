package common

import (
	"fmt"
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
	Position [3]float32
	Normal   [3]float32
	Tint     [4]uint8
	Uv       [2]float32
	Uv2      [2]float32
	Bone     [4]uint16
	Weight   [4]float32
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
