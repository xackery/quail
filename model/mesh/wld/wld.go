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
	archive        archive.ReadWriter
	name           string
	BspRegionCount uint32
	Hash           map[int]string
	materials      []*geo.Material
	meshes         []*mesh
	NameCache      map[int32]string
}

type mesh struct {
	name      string
	vertices  []*geo.Vertex
	triangles []*geo.Triangle
}

type fragmentInfo struct {
	name string
	data archive.WldFragmenter
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.ReadWriter) (*WLD, error) {
	e := &WLD{
		name:    name,
		archive: pfs,
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*WLD, error) {
	e := &WLD{
		name:    name,
		archive: pfs,
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

func (e *WLD) SetLayers(layers []*geo.Layer) error {
	return nil
}

func (e *WLD) SetParticleRenders(particles []*geo.ParticleRender) error {
	return nil
}

// SetParticlePoints sets particle points for a world file
func (e *WLD) SetParticlePoints(particles []*geo.ParticlePoint) error {
	return nil
}
