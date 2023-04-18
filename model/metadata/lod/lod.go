// lod is a package that contains the LOD file details
package lod

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs/archive"
)

// LOD is level of detail information
// Typical usaeg is like so:
/*
EQLOD
LOD,OBJ_FIREPIT_STMFT,150
LOD,OBJ_FIREPIT_STMFT_LOD1,250
LOD,OBJ_FIREPIT_STMFT_LOD2,400
LOD,OBJ_FIREPIT_STMFT_LOD3,1000
*/
type LOD struct {
	name string
	pfs  archive.Reader
	lods []*LODEntry
}

type LODEntry struct {
	Category   string
	ObjectName string
	Distance   int
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.Reader) (*LOD, error) {
	t := &LOD{
		name: name,
	}
	return t, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.Reader, file string) (*LOD, error) {
	e := &LOD{
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
func (e *LOD) Name() string {
	return e.name
}

// Data returns the raw data of the file
func (e *LOD) Data() []byte {
	w := bytes.NewBuffer(nil)

	err := e.Encode(w)
	if err != nil {
		log.Errorf("Failed to encode litrain data: %s", err)
		os.Exit(1)
	}
	return w.Bytes()
}
