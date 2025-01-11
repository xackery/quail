package common

import (
	"io"
)

// FragmentReadWriter is used to read a fragment in wld format
type FragmentReadWriter interface {
	FragmentReader
	FragmentWriter
	TreeLinker
}

type FragmentReader interface {
	Read(w io.ReadSeeker, isNewWorld bool) error
	FragCode() int
}

// FragmentWriter2 is used to write a fragment in wld format
type FragmentWriter interface {
	Write(w io.Writer, isNewWorld bool) error
}
