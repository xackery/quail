package qfs

import (
	"io"
	"io/fs"
)

// QFS is an abstraction over a filesystem.
type QFS interface {
	fs.FS                                       // Provides Open(name string) (fs.File, error)
	Stat(name string) (fs.FileInfo, error)      // Provides file metadata
	ReadDir(name string) ([]fs.DirEntry, error) // Lists directory contents
	ReadFile(name string) ([]byte, error)       // Reads a file
	RemoveAll(name string) error                // Removes a file or directory
	MkdirAll(name string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	Create(name string) (io.WriteCloser, error)
}
