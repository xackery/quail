package common

import "io"

type GLTFExporter interface {
	GLTFExport(w io.Writer) error
}
