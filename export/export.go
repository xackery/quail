// export parses eqg files and smartly prepares models for export
package export

import (
	"io"

	"github.com/xackery/quail/common"
	qgltf "github.com/xackery/quail/gltf"
)

type Export struct {
	name    string
	archive common.ArchiveReadWriter
	model   modeler
}

type modeler interface {
	Load(r io.ReadSeeker) error
	GLTFExport(doc *qgltf.GLTF) error
	SetLayers(layers []*common.Layer) error
	SetParticleRenders(particles []*common.ParticleRender) error
	SetParticlePoints(particles []*common.ParticlePoint) error
}

func New(name string, archive common.ArchiveReadWriter) (*Export, error) {
	return &Export{
		name:    name,
		archive: archive,
	}, nil
}
