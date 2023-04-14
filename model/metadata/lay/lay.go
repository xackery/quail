// lay is an EverQuest layer file, giving details about layered models
package lay

// https://github.com/Zaela/EQGModelImporter/blob/master/src/mds.cpp

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// LAY is a layer definition
type LAY struct {
	// name is used as an identifier
	name    string
	version uint32
	// pfs is used as an alternative to path when loading data from a pfs file
	pfs          archive.Reader
	layerManager *geo.LayerManager
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.Reader) (*LAY, error) {
	e := &LAY{
		name:         name,
		pfs:          pfs,
		layerManager: &geo.LayerManager{},
	}
	return e, nil
}

func NewFile(name string, pfs archive.Reader, file string) (*LAY, error) {
	e := &LAY{
		name:         name,
		pfs:          pfs,
		layerManager: &geo.LayerManager{},
	}
	data, err := pfs.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", file, err)
	}
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return e, nil
}

func (e *LAY) SetName(value string) {
	e.name = value
}

// Layers returns a list of layers
func (e *LAY) Layers() []*geo.Layer {
	return e.layerManager.Layers()
}

// LayerByIndex returns a layer by index
func (e *LAY) LayerByIndex(index int) *geo.Layer {
	layers := e.layerManager.Layers()
	if len(layers) <= index {
		return nil
	}
	return layers[index]
}

// LayerCount returns the number of layers
func (e *LAY) LayerCount() int {
	return e.layerManager.Count()
}

// Name returns the name of the file
func (e *LAY) Name() string {
	return e.name
}
