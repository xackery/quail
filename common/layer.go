package common

import "fmt"

type Layer struct {
	Name   string
	Entry0 string
	Entry1 string
}

func (e *Layer) String() string {
	return fmt.Sprintf("%s %s %s", e.Name, e.Entry0, e.Entry1)
}
