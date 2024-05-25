package model

import "io"

// FragmentReadWriter is used to read a fragment in wld format
type FragmentReadWriter interface {
	FragmentReader
	FragmentWriter
}

type FragmentReader interface {
	Read(w io.ReadSeeker) error
	FragCode() int
}

// FragmentWriter2 is used to write a fragment in wld format
type FragmentWriter interface {
	Write(w io.Writer) error
}
