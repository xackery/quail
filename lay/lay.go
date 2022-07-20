// lay is an EverQuest layer file, giving details about layered models
package lay

// https://github.com/Zaela/EQGModelImporter/blob/master/src/mds.cpp

import (
	"github.com/xackery/quail/common"
)

// LAY is a zon file struct
type LAY struct {
	// name is used as an identifier
	name string
	// archive is used as an alternative to path when loading data from a archive file
	archive common.ArchiveReader
	layers  []*common.Layer
}

func New(name string, archive common.ArchiveReader) (*LAY, error) {
	e := &LAY{
		name:    name,
		archive: archive,
	}
	return e, nil
}

func (e *LAY) SetName(value string) {
	e.name = value
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
