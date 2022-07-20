// eqg is a pfs archive for EverQuest
package eqg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/common"
)

// EQG represents a modern everquest zone archive format
type EQG struct {
	name      string
	files     []common.Filer
	fileCount int
}

// New creates an empty eqg archive
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
	err = e.Load(r)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	return e, nil
}
