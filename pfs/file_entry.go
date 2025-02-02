package pfs

import "path/filepath"

// FileEntry represents a file entry in a Pfs
type FileEntry struct {
	name string
	data []byte
}

// NewFileEntry creates a new file entry
func NewFileEntry(name string, data []byte) *FileEntry {
	return &FileEntry{
		name: name,
		data: data,
	}
}

// SetName sets the name of the file entry
func (e *FileEntry) SetName(name string) error {
	e.name = name
	return nil
}

// Name returns the name of the file entry
func (e *FileEntry) Name() string {
	return e.name
}

// Ext returns the extension of a file entry
func (e *FileEntry) Ext() string {
	return filepath.Ext(e.name)
}

// SetData sets the data of the file entry
func (e *FileEntry) SetData(data []byte) error {
	e.data = data
	return nil
}

// Data returns the data of the file entry
func (e *FileEntry) Data() []byte {
	return e.data
}
