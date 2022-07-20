package common

// ArchiveReadWriter implements both reader and writer operations to a writer
type ArchiveReadWriter interface {
	ArchiveReader
	ArchiveWriter
}

// ArchiveReader is used to access data from an archive (or a flat file)
type ArchiveReader interface {
	File(name string) ([]byte, error)
	Files() []Filer
	Len() int
}

// ArchiveWriter is used to write data to an archive (or flat file)
type ArchiveWriter interface {
	WriteFile(name string, data []byte) error
}
