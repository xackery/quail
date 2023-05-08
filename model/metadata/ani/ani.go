// ani are animation files, found in EverQuest eqg files
package ani

import (
	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

type ANI struct {
	name     string
	bones    []*geo.BoneAnimation
	isStrict bool
	pfs      archive.ReadWriter
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.ReadWriter) (*ANI, error) {
	e := &ANI{
		name: name,
		pfs:  pfs,
	}
	return e, nil
}

// Name returns the name of the model
func (e *ANI) Name() string {
	return e.name
}
