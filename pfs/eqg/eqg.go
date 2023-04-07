// eqg is a pfs archive for EverQuest
package eqg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/pfs/archive"
)

// EQG represents a modern everquest zone archive format
type EQG struct {
	name      string
	files     []archive.Filer
	fileCount int
}

// New creates a new empty instance. Use NewFile to load an archive on creation
func New(name string) (*EQG, error) {
	e := &EQG{
		name: name,
	}
	return e, nil
}

// NewFile takes path and loads it as an eqg archive
func NewFile(path string) (*EQG, error) {
	e := &EQG{
		name: filepath.Base(path),
	}
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	err = e.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return e, nil
}
