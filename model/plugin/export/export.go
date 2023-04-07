// export parses eqg files and smartly prepares models for export
package export

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

type Export struct {
	name  string
	pfs   archive.ReadWriter
	model modeler
}

type modeler interface {
	Decode(r io.ReadSeeker) error
	BlenderExport(dir string) error
	SetLayers(layers []*geo.Layer) error
	SetParticleRenders(particles []*geo.ParticleRender) error
	SetParticlePoints(particles []*geo.ParticlePoint) error
}

// New creates a new empty instance. Use NewFile to load the archive on creation
func New(name string, pfs archive.ReadWriter) (*Export, error) {
	return &Export{
		name: name,
		pfs:  pfs,
	}, nil
}

func NewFile(path string, pfs archive.ReadWriter) (*Export, error) {
	e := &Export{
		name: filepath.Base(path),
		pfs:  pfs,
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
