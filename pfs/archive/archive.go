package archive

import "io"

// ReadWriter implements both reader and writer operations to a writer
type ReadWriter interface {
	Reader
	Writer
}

// Reader is used to access data from an archive (or a flat file)
type Reader interface {
	File(name string) ([]byte, error)
	Files() []Filer
	Len() int
}

// Writer is used to write data to an archive (or flat file)
type Writer interface {
	WriteFile(name string, data []byte) error
	Encode(w io.WriteSeeker) error
}
