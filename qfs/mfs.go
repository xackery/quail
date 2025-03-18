package qfs

import (
	"io"
	"io/fs"
	"strings"
	"time"
)

// MFS is an in-memory filesystem.
type MFS struct {
	files map[string]*memFile
}

type memFile struct {
	data []byte
	info fileInfo
}

type fileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
}

func (m *fileInfo) Name() string               { return m.name }
func (m *fileInfo) Size() int64                { return m.size }
func (m *fileInfo) Mode() fs.FileMode          { return m.mode }
func (m *fileInfo) ModTime() time.Time         { return m.modTime }
func (m *fileInfo) IsDir() bool                { return m.mode.IsDir() }
func (m *fileInfo) Sys() any                   { return nil }
func (m *fileInfo) Type() fs.FileMode          { return m.mode }
func (m *fileInfo) Info() (fs.FileInfo, error) { return m, nil }

func NewMemoryFS() *MFS {
	return &MFS{files: make(map[string]*memFile)}
}

// Open implements fs.FS and returns a file reader.
func (m *MFS) Open(name string) (fs.File, error) {
	f, ok := m.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return &memReadWriter{data: f.data, info: &f.info}, nil
}

// Stat returns file metadata.
func (m *MFS) Stat(name string) (fs.FileInfo, error) {
	f, ok := m.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return &f.info, nil
}

// ReadDir lists directory contents.
func (m *MFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var entries []fs.DirEntry
	for path, f := range m.files {
		if strings.HasPrefix(path, name+"/") {
			entries = append(entries, &fileInfo{name: f.info.name, mode: f.info.mode})
		}
	}
	return entries, nil
}

// ReadFile reads a file.
func (m *MFS) ReadFile(name string) ([]byte, error) {
	f, ok := m.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return f.data, nil
}

func (m *MFS) RemoveAll(name string) error {
	delete(m.files, name)
	return nil
}

func (m *MFS) MkdirAll(name string, perm fs.FileMode) error {
	return nil
}

func (m *MFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	m.files[name] = &memFile{data: data, info: fileInfo{name: name, size: int64(len(data)), mode: perm, modTime: time.Now()}}
	return nil
}

func (m *MFS) Create(name string) (io.WriteCloser, error) {
	f := &memFile{info: fileInfo{name: name, mode: 0666, modTime: time.Now()}}
	m.files[name] = f
	return &memReadWriter{data: f.data, info: &f.info}, nil
}

// memReadWriter represents an in-memory file reader.
type memReadWriter struct {
	data  []byte
	index int
	info  *fileInfo
}

// Read implements io.Reader.
func (m *memReadWriter) Read(p []byte) (int, error) {
	if m.index >= len(m.data) {
		return 0, io.EOF
	}
	n := copy(p, m.data[m.index:])
	m.index += n
	return n, nil
}

// Close implements io.Closer.
func (m *memReadWriter) Close() error { return nil }

// Stat implements fs.File by returning file info.
func (m *memReadWriter) Stat() (fs.FileInfo, error) {
	return m.info, nil
}

// Write implements io.Writer.
func (m *memReadWriter) Write(p []byte) (int, error) {
	m.data = append(m.data, p...)
	return len(p), nil
}
