package qfs

import (
	"io"
	"io/fs"
	"os"
)

// OSFS is an implementation of FS using the local filesystem.
type OSFS struct{}

func (OSFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (OSFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (OSFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (OSFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (OSFS) RemoveAll(name string) error {
	return os.RemoveAll(name)
}

func (OSFS) MkdirAll(name string, perm fs.FileMode) error {
	return os.MkdirAll(name, perm)
}

func (OSFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (OSFS) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}
