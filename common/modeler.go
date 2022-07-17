package common

type Modeler interface {
	SetLayers(layers []*Layer) error
	AddFile(fe *FileEntry)
	MaterialAdd(name string, shaderName string) error
	MaterialPropertyAdd(materialName string, propertyName string, category uint32, value string) error
}
