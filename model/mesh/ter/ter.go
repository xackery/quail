// ter is an EverQuest terrain model file
package ter

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// TER is a terrain file struct
type TER struct {
	name            string
	version         uint32
	meshManager     *geo.MeshManager
	MaterialManager *geo.MaterialManager
	particleManager *geo.ParticleManager
	files           []archive.Filer
	archive         archive.ReadWriter
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.ReadWriter) (*TER, error) {
	t := &TER{
		name:            name,
		archive:         pfs,
		MaterialManager: geo.NewMaterialManager(),
		meshManager:     &geo.MeshManager{},
		particleManager: &geo.ParticleManager{},
	}
	return t, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*TER, error) {
	e := &TER{
		name:            name,
		archive:         pfs,
		MaterialManager: geo.NewMaterialManager(),
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

func (e *TER) Name() string {
	return e.name
}

func (e *TER) Data() []byte {
	w := bytes.NewBuffer(nil)
	err := e.Encode(w)
	if err != nil {
		log.Warnf("Failed to encode terrain data: %s", err)
		os.Exit(1)
	}
	return w.Bytes()
}

func (e *TER) SetName(value string) {
	e.name = value
}

// Close flushes the data in a mod
func (e *TER) Close() {
	e.files = nil
	e.MaterialManager = &geo.MaterialManager{}
	e.meshManager = &geo.MeshManager{}
	e.particleManager = &geo.ParticleManager{}
}
