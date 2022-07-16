package s3d

import (
	"fmt"

	"github.com/xackery/quail/common"
)

// File returns data of a file
func (e *S3D) File(name string) ([]byte, error) {
	for _, f := range e.files {
		if f.Name() == name {
			return f.Data(), nil
		}
	}
	return nil, fmt.Errorf("%s not found", name)
}

// Files returns a string array of every file inside an EQG
func (e *S3D) Files() []common.Filer {
	return e.files
}
