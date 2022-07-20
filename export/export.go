// export parses eqg files and smartly prepares models for export
package export

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/xackery/quail/common"
	qgltf "github.com/xackery/quail/gltf"
)

type Export struct {
	name    string
	archive common.ArchiveReadWriter
	model   modeler
}

type modeler interface {
	Decode(r io.ReadSeeker) error
	GLTFEncode(doc *qgltf.GLTF) error
	SetLayers(layers []*common.Layer) error
	SetParticleRenders(particles []*common.ParticleRender) error
	SetParticlePoints(particles []*common.ParticlePoint) error
}

// New creates a new empty instance. Use NewFile to load the archive on creation
func New(name string, archive common.ArchiveReadWriter) (*Export, error) {
	return &Export{
		name:    name,
		archive: archive,
	}, nil
}

func NewFile(path string, archive common.ArchiveReadWriter) (*Export, error) {
	e := &Export{
		name:    filepath.Base(path),
		archive: archive,
	}

	/*data, err := archive.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %s", file, err)
	}*/
	err := e.LoadArchive()
	if err != nil {
		return nil, fmt.Errorf("loadArchive: %w", err)
	}

	return e, nil
}

func (e *Export) Name() string {
	return e.name
}
