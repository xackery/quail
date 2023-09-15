package s3d

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/pfs/archive"
)

// File returns data of a file
func (e *S3D) File(name string) ([]byte, error) {
	name = strings.ToLower(name)
	for _, f := range e.files {
		if strings.EqualFold(f.Name(), name) {
			return f.Data(), nil
		}
	}
	return nil, fmt.Errorf("%s not found", name)
}

// Files returns a string array of every file inside an EQG
func (e *S3D) Files() []archive.Filer {
	return e.files
}

func (e *S3D) Len() int {
	return e.fileCount
}

func (e *S3D) WriteFile(name string, data []byte) error {
	name = strings.ToLower(name)
	for _, file := range e.files {
		if strings.EqualFold(file.Name(), name) {
			return file.SetData(data)
		}
	}
	fe, err := archive.NewFileEntry(name, data)
	if err != nil {
		return fmt.Errorf("newFileEntry: %w", err)
	}
	e.files = append(e.files, fe)
	return nil
}
