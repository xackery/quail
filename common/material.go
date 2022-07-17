package common

import (
	"fmt"
	"strings"
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

type MaterialByName []*Material

func (s MaterialByName) Len() int {
	return len(s)
}

func (s MaterialByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s MaterialByName) Less(i, j int) bool {
	return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name)
}
