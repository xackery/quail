// mod is an EverQuest model file format
package mod

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// MOD is a zon file struct
type MOD struct {
	// name is used as an identifier
	name string
	// path is used for relative paths when looking for flat file texture references
	path string
	// pfs is used as an alternative to path when loading data from a pfs file
	pfs             archive.Reader
	MaterialManager *geo.MaterialManager
	meshManager     *geo.MeshManager
	particleManager *geo.ParticleManager
	files           []archive.Filer // list of files known to be linked in this mod
	version         uint32
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.Reader) (*MOD, error) {
	e := &MOD{
		name: name,
		pfs:  pfs,
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*MOD, error) {
	e := &MOD{
		name: name,
		pfs:  pfs,
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

func (e *MOD) SetName(value string) {
	e.name = value
}

func (e *MOD) SetPath(value string) {
	e.path = value
}

func (e *MOD) SetLayers(layers []*geo.Layer) error {
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
			err = e.MaterialManager.PropertyAdd(o.Name, "e_texturediffuse0", 2, normalName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", normalName, err)
			}
		}
	}
	return nil
}

func (e *MOD) AddFile(fe *archive.FileEntry) {
	e.files = append(e.files, fe)
}

func (e *MOD) Name() string {
	return e.name
}

// Close flushes the data in a mod
func (e *MOD) Close() {
	e.files = nil
	e.MaterialManager = &geo.MaterialManager{}
	e.meshManager = &geo.MeshManager{}
	e.particleManager = &geo.ParticleManager{}
}
