package mod

import (
	"fmt"
	"strings"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

// MOD is a zon file struct
type MOD struct {
	// name is used as an identifier
	name string
	// path is used for relative paths when looking for flat file texture references
	path string
	// eqg is used as an alternative to path when loading data from a eqg file
	eqg            common.Archiver
	materials      []*common.Material
	vertices       []*common.Vertex
	faces          []*common.Face
	bones          []*bone
	files          []common.Filer
	gltfBoneBuffer map[int]uint32
}

type bone struct {
	name          string
	next          int32
	childrenCount uint32
	childIndex    int32
	pivot         *math32.Vector3
	rot           *math32.Vector4
	scale         *math32.Vector3
}

func New(name string, path string) (*MOD, error) {
	e := &MOD{
		name: name,
		path: path,
	}
	return e, nil
}

func NewEQG(name string, eqg common.Archiver) (*MOD, error) {
	e := &MOD{
		name: name,
		eqg:  eqg,
	}
	return e, nil
}

func (e *MOD) SetName(value string) {
	e.name = value
}

func (e *MOD) SetPath(value string) {
	e.path = value
}

func (e *MOD) SetLayers(layers []*common.Layer) error {
	for _, o := range layers {
		err := e.MaterialAdd(o.Name, "")
		if err != nil {
			return fmt.Errorf("materialAdd: %w", err)
		}
		entry0Name := strings.ToLower(o.Entry0)
		entry1Name := strings.ToLower(o.Entry1)
		diffuseName := ""
		normalName := ""
		if strings.Contains(entry0Name, "_c.dds") {
			diffuseName = entry0Name
		}
		if strings.Contains(entry1Name, "_c.dds") {
			diffuseName = entry1Name
		}

		if strings.Contains(entry0Name, "_n.dds") {
			normalName = entry0Name
		}
		if strings.Contains(entry1Name, "_n.dds") {
			normalName = entry1Name
		}

		if len(diffuseName) > 0 {
			err = e.MaterialPropertyAdd(o.Name, "e_texturediffuse0", 2, diffuseName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", diffuseName, err)
			}
		}

		if len(normalName) > 0 {
			err = e.MaterialPropertyAdd(o.Name, "e_texturediffuse0", 2, normalName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", normalName, err)
			}
		}
	}
	return nil
}

func (e *MOD) AddFile(fe *common.FileEntry) {
	e.files = append(e.files, fe)
}
