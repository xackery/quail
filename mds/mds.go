// mds is an EverQuest model file format
package mds

// https://github.com/Zaela/EQGModelImporter/blob/master/src/mds.cpp

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

// MDS is a zon file struct
type MDS struct {
	// base is the mds's base model name
	name string
	// path is used for relative paths when looking for flat file texture references
	path string
	// archive is used as an alternative to path when loading data from a archive file
	archive         common.ArchiveReadWriter
	materials       []*common.Material
	vertices        []*common.Vertex
	faces           []*common.Face
	bones           []*bone
	files           []common.Filer
	joints          [][4]uint16
	weights         [][4]float32
	particleRenders []*common.ParticleRender
	particlePoints  []*common.ParticlePoint
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

func New(name string, archive common.ArchiveReadWriter) (*MDS, error) {
	e := &MDS{
		name:    name,
		archive: archive,
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, archive common.ArchiveReadWriter, file string) (*MDS, error) {
	e := &MDS{
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

func (e *MDS) SetName(value string) {
	e.name = value
}

func (e *MDS) SetPath(value string) {
	e.path = value
}

func (e *MDS) SetLayers(layers []*common.Layer) error {
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
			err = e.MaterialPropertyAdd(o.Name, "e_texturenormal0", 2, normalName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", normalName, err)
			}
		}
	}
	sort.Sort(common.MaterialByName(e.materials))
	return nil
}

func (e *MDS) AddFile(fe *common.FileEntry) {
	e.files = append(e.files, fe)
}

func (e *MDS) SetParticleRenders(particles []*common.ParticleRender) error {
	e.particleRenders = particles
	return nil
}

func (e *MDS) SetParticlePoints(particles []*common.ParticlePoint) error {
	e.particlePoints = particles
	return nil
}

func (e *MDS) Name() string {
	return e.name
}
