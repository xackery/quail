package geo

import (
	"github.com/xackery/quail/pfs/archive"
)

// Modeler is a modeler interface
type Modeler interface {
	SetLayers(layers []*Layer) error
	AddFile(fe *archive.FileEntry)
	MaterialAdd(name string, shaderName string) error
	MaterialPropertyAdd(materialName string, propertyName string, category uint32, value string) error
}
