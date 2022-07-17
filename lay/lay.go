package lay

// https://github.com/Zaela/EQGModelImporter/blob/master/src/mds.cpp

import (
	"github.com/xackery/quail/common"
)

// LAY is a zon file struct
type LAY struct {
	// name is used as an identifier
	name string
	// path is used for relative paths when looking for flat file texture references
	path string
	// eqg is used as an alternative to path when loading data from a eqg file
	eqg    common.Archiver
	layers []*common.Layer
}

func New(name string, path string) (*LAY, error) {
	e := &LAY{
		name: name,
		path: path,
	}
	return e, nil
}

func NewEQG(name string, eqg common.Archiver) (*LAY, error) {
	e := &LAY{
		name: name,
		eqg:  eqg,
	}
	return e, nil
}

func (e *LAY) SetName(value string) {
	e.name = value
}

func (e *LAY) SetPath(value string) {
	e.path = value
}

func (e *LAY) Layers() []*common.Layer {
	return e.layers
}

func (e *LAY) LayerByIndex(index int) *common.Layer {
	if len(e.layers) <= index {
		return nil
	}
	return e.layers[index]
}

func (e *LAY) LayerCount() int {
	return len(e.layers)
}
