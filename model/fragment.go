package model

import "io"

// FragmentReadWriter is used to read a fragment in wld format
type FragmentReadWriter interface {
	FragmentReader
	FragmentWriter
}

type FragmentReader interface {
	Read(w io.ReadSeeker, isNewWorld bool) error
	FragCode() int
	NameRef() int32
}

// FragmentWriter2 is used to write a fragment in wld format
type FragmentWriter interface {
	Write(w io.Writer, isNewWorld bool) error
}
