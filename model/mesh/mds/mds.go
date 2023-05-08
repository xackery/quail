// mds is an EverQuest model file format
package mds

// https://github.com/Zaela/EQGModelImporter/blob/master/src/mds.cpp

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// MDS is a zon file struct
type MDS struct {
	// base is the mds's base model name
	name string
	// path is used for relative paths when looking for flat file texture references
	path            string
	itemName        string // Contains an origin item, seems optional, unsure where it's reffered to
	version         uint32
	pfs             archive.ReadWriter
	files           []archive.Filer
	MaterialManager *geo.MaterialManager
	meshManager     *geo.MeshManager
	particleManager *geo.ParticleManager
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.ReadWriter) (*MDS, error) {
	e := &MDS{
		name:            name,
		pfs:             pfs,
		MaterialManager: geo.NewMaterialManager(),
		meshManager:     geo.NewMeshManager(),
		particleManager: geo.NewParticleManager(),
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*MDS, error) {
	e := &MDS{
		name:            name,
		pfs:             pfs,
		MaterialManager: geo.NewMaterialManager(),
		meshManager:     geo.NewMeshManager(),
		particleManager: geo.NewParticleManager(),
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

func (e *MDS) SetName(value string) {
	e.name = value
}

func (e *MDS) SetPath(value string) {
	e.path = value
}

func (e *MDS) AddFile(fe *archive.FileEntry) {
	e.files = append(e.files, fe)
}

func (e *MDS) Name() string {
	return e.name
}

// Close flushes the data in a mod
func (e *MDS) Close() {
	e.files = nil
	e.MaterialManager = geo.NewMaterialManager()
	e.meshManager = geo.NewMeshManager()
	e.particleManager = geo.NewParticleManager()
}
