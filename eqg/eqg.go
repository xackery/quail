package eqg

import (
	"github.com/xackery/quail/common"
)

// EQG represents a modern everquest zone archive format
type EQG struct {
	name  string
	files []common.Filer
}

func New(name string) (*EQG, error) {
	e := &EQG{
		name: name,
	}
	return e, nil
}
