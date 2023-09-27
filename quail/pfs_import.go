package quail

import (
	"fmt"
	"path/filepath"
)

// Import imports the quail target
func (e *Quail) PFSImport(path string) error {
	ext := filepath.Ext(path)

	switch ext {
	case ".eqg":
		return e.EQGImport(path)
	case ".s3d":
		return e.S3DImport(path)
	case ".pfs":
		return e.EQGImport(path)
	case ".pak":
		return e.EQGImport(path)
	default:
		return fmt.Errorf("unknown pfs type %s, valid options are eqg and pfs", ext)
	}
}
