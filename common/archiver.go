package common

// Archiver is used to access data from an archive
type Archiver interface {
	File(name string) ([]byte, error)
}
