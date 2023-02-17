// wld contains EverQuest fragments for various data
package wld

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
)

// WLD is a wld file struct
type WLD struct {
	archive         common.ArchiveReadWriter
	name            string
	BspRegionCount  uint32
	Hash            map[int]string
	materials       []*common.Material
	files           []common.Filer
	particleRenders []*common.ParticleRender
	particlePoints  []*common.ParticlePoint
	meshes          []*mesh
	NameCache       map[int32]string
}

type mesh struct {
	name      string
	vertices  []*common.Vertex
	triangles []*common.Triangle
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

func (e *WLD) SetLayers(layers []*common.Layer) error {
	return nil
}

func (e *WLD) SetParticleRenders(particles []*common.ParticleRender) error {
	return nil
}

// SetParticlePoints sets particle points for a world file
func (e *WLD) SetParticlePoints(particles []*common.ParticlePoint) error {
	return nil
}
