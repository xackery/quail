// s3d is an EverQuest pfs archive
package s3d

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/common"
)

// S3D represents a classic everquest zone archive format
type S3D struct {
	name      string
	ShortName string
	files     []common.Filer
	fileCount int
}

// New creates a new empty instance. Use NewFile to load an archive on creation
func New(name string) (*S3D, error) {
	e := &S3D{
		name: name,
	}
	return e, nil
}

// NewFile takes path and loads it as an eqg archive
func NewFile(path string) (*S3D, error) {
	e := &S3D{
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
