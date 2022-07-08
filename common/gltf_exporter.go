package common

import "io"

type GLTFExporter interface {
	ExportGLTF(w io.Writer) error
}
