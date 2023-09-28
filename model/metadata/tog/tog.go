// tog contains terrain details
package tog

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/pfs"
)

type TOG struct {
	name    string
	pfs     *pfs.PFS
	objects []common.Object
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs *pfs.PFS) (*TOG, error) {
	e := &TOG{
		name: name,
		pfs:  pfs,
	}
	return e, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs *pfs.PFS, file string) (*TOG, error) {
	e := &TOG{
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

// Name returns the name of the file
func (e *TOG) Name() string {
	return e.name
}
