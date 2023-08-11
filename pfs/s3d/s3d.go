// s3d is an EverQuest pfs archive
package s3d

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xackery/quail/pfs/archive"
)

// S3D represents a classic everquest zone archive format
type S3D struct {
	name            string
	ShortName       string
	files           []archive.Filer
	fileCount       int
	ContentsSummary string
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

func (e *S3D) Close() error {
	e.ContentsSummary = "\n"
	for i, fe := range e.files {
		base := float64(len(fe.Data()))
		out := ""
		num := float64(1024)
		if base < num*num*num*num {
			out = fmt.Sprintf("%0.0fG", base/num/num/num)
		}
		if base < num*num*num {
			out = fmt.Sprintf("%0.0fM", base/num/num)
		}
		if base < num*num {
			out = fmt.Sprintf("%0.0fK", base/num)
		}
		if base < num {
			out = fmt.Sprintf("%0.0fB", base)
		}
		e.ContentsSummary += fmt.Sprintf("%d %s:\t %s\n", i, out, fe.Name())
	}
	e.files = nil
	e.name = ""
	e.fileCount = 0
	return nil
}
