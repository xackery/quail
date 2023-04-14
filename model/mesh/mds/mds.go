// mds is an EverQuest model file format
package mds

// https://github.com/Zaela/EQGModelImporter/blob/master/src/mds.cpp

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// MDS is a zon file struct
type MDS struct {
	// base is the mds's base model name
	name string
	// path is used for relative paths when looking for flat file texture references
	path string
	// pfs is used as an alternative to path when loading data from a pfs file
	version         uint32
	pfs             archive.ReadWriter
	files           []archive.Filer
	isDecoded       bool
	MaterialManager *geo.MaterialManager
	meshManager     *geo.MeshManager
	particleManager *geo.ParticleManager
	animations      []*geo.BoneAnimation
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.ReadWriter) (*MDS, error) {
	e := &MDS{
		name:            name,
		pfs:             pfs,
		MaterialManager: &geo.MaterialManager{},
		meshManager:     &geo.MeshManager{},
		particleManager: &geo.ParticleManager{},
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*MDS, error) {
	e := &MDS{
		name:            name,
		pfs:             pfs,
		MaterialManager: &geo.MaterialManager{},
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

func (e *MDS) SetName(value string) {
	e.name = value
}

func (e *MDS) SetPath(value string) {
	e.path = value
}

func (e *MDS) SetLayers(layers []*geo.Layer) error {
	for _, o := range layers {
		err := e.MaterialManager.Add(o.Name, "")
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
			err = e.MaterialManager.PropertyAdd(o.Name, "e_texturediffuse0", 2, diffuseName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", diffuseName, err)
			}
		}

		if len(normalName) > 0 {
			err = e.MaterialManager.PropertyAdd(o.Name, "e_texturenormal0", 2, normalName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", normalName, err)
			}
		}
	}
	e.MaterialManager.SortByName()
	return nil
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
	e.MaterialManager = &geo.MaterialManager{}
	e.meshManager = &geo.MeshManager{}
	e.particleManager = &geo.ParticleManager{}
}
