// wld contains EverQuest fragments for various data
package wld

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// WLD is a wld file struct
type WLD struct {
	archive         archive.ReadWriter
	version         uint32
	name            string
	meshManager     *geo.MeshManager
	materialManager *geo.MaterialManager
	particleManager *geo.ParticleManager
	BspRegionCount  uint32
	Hash            map[int]string
	NameCache       map[int32]string
}

type fragmentInfo struct {
	name string
	data archive.WldFragmenter
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.ReadWriter) (*WLD, error) {
	e := &WLD{
		name:            name,
		archive:         pfs,
		materialManager: &geo.MaterialManager{},
		meshManager:     &geo.MeshManager{},
		particleManager: &geo.ParticleManager{},
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*WLD, error) {
	e := &WLD{
		name:            name,
		archive:         pfs,
		materialManager: &geo.MaterialManager{},
		meshManager:     &geo.MeshManager{},
		particleManager: &geo.ParticleManager{},
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

// SetLayers sets layers for a world file
func (e *WLD) SetLayers(layers []*geo.Layer) error {
	return nil
}

// SetParticleRenders sets particle renders for a world file
func (e *WLD) SetParticleRenders(particles []*geo.ParticleRender) error {
	return nil
}

// SetParticlePoints sets particle points for a world file
func (e *WLD) SetParticlePoints(particles []*geo.ParticlePoint) error {
	return nil
}

// Name returns the name of the archive
func (e *WLD) Name() string {
	return e.name
}
