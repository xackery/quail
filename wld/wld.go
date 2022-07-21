// wld contains EverQuest fragments for various data
package wld

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
)

// WLD is a wld file struct
type WLD struct {
	archive            common.ArchiveReadWriter
	name               string
	BspRegionCount     uint32
	Hash               map[int]string
	fragments          []*fragmentInfo
	materials          []*common.Material
	vertices           []*common.Vertex
	triangles          []*common.Triangle
	files              []common.Filer
	gltfMaterialBuffer map[string]*uint32
	gltfBoneBuffer     map[int]uint32
}

type fragmentInfo struct {
	name string
	data common.WldFragmenter
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, archive common.ArchiveReadWriter) (*WLD, error) {
	e := &WLD{
		name:    name,
		archive: archive,
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, archive common.ArchiveReadWriter, file string) (*WLD, error) {
	e := &WLD{
		name:    name,
		archive: archive,
	}
	data, err := archive.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", file, err)
	}
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return e, nil
}
