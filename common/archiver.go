package common

// Archiver is used to access data from an archive (or a flat file)
type Archiver interface {
	File(name string) ([]byte, error)
	Files() []Filer
	Len() int
}
