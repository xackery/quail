//go:build tinygo.wasm

package os

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"syscall/js"
	"time"

	nativeos "os"
)

var ErrNotExist = nativeos.ErrNotExist
var ModePerm = nativeos.ModePerm

type File struct {
	Name    string
	Content []byte
	Perm    fs.FileMode
	ModTime time.Time
	IsDir   bool
	Entries map[string]*File
	cursor  int
}

var fileSystem = map[string]*File{
	"/": {Name: "/", IsDir: true, Entries: make(map[string]*File)},
}

func ExportFileSystem() js.Value {
	jsFileSystem := js.Global().Get("Object").New()

	var exportDirectory func(currentPath string, dir *File, jsDir js.Value)
	exportDirectory = func(currentPath string, dir *File, jsDir js.Value) {
		for name, file := range dir.Entries {
			fullPath := currentPath + "/" + name
			if file.IsDir {
				subDir := js.Global().Get("Object").New()
				jsDir.Set(name, subDir)
				exportDirectory(fullPath, file, subDir)
			} else {
				jsDir.Set(name, string(file.Content))
			}
		}
	}

	exportDirectory("", fileSystem["/"], jsFileSystem)
	return jsFileSystem
}

func (f *File) Read(b []byte) (int, error) {
	if f.IsDir {
		return 0, errors.New("cannot read from a directory")
	}
	if f.Content == nil || f.cursor >= len(f.Content) {
		return 0, io.EOF
	}
	n := copy(b, f.Content[f.cursor:])
	f.cursor += n
	return n, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.IsDir {
		return 0, errors.New("cannot seek in a directory")
	}
	var newCursor int
	switch whence {
	case io.SeekStart:
		newCursor = int(offset)
	case io.SeekCurrent:
		newCursor = f.cursor + int(offset)
	case io.SeekEnd:
		newCursor = len(f.Content) + int(offset)
	default:
		return 0, errors.New("invalid whence")
	}
	if newCursor < 0 || newCursor > len(f.Content) {
		return 0, errors.New("invalid offset")
	}
	f.cursor = newCursor
	return int64(f.cursor), nil
}

func (f *File) Write(b []byte) (int, error) {
	if f.IsDir {
		return 0, errors.New("cannot write to a directory")
	}
	if f.cursor+len(b) > len(f.Content) {
		newContent := make([]byte, f.cursor+len(b))
		copy(newContent, f.Content)
		f.Content = newContent
	}
	n := copy(f.Content[f.cursor:], b)
	f.cursor += n
	f.ModTime = time.Now()
	return n, nil
}

func (f *File) Close() error {
	return nil
}

func (f *File) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}

func Stat(name string) (fs.FileInfo, error) {
	file, exists := fileSystem[name]
	if !exists {
		return nil, errors.New("file does not exist")
	}
	return fileInfo{file}, nil
}

func ReadDir(name string) ([]fs.DirEntry, error) {
	dir, exists := fileSystem[name]
	if !exists || !dir.IsDir {
		return nil, errors.New("read directory does not exist " + name)
	}
	entries := make([]fs.DirEntry, 0, len(dir.Entries))
	for _, entry := range dir.Entries {
		entries = append(entries, dirEntry{entry})
	}
	return entries, nil
}

func WriteFile(name string, buffer []byte, perm fs.FileMode) error {
	file, exists := fileSystem[name]
	if !exists {
		dirPath := filepath.Dir(name)
		dir, exists := fileSystem[dirPath]
		if !exists || !dir.IsDir {
			return errors.New("directory does not exist " + name + " dir path " + dirPath)
		}
		file = &File{Name: name, Perm: perm, ModTime: time.Now()}
		dir.Entries[filepath.Base(name)] = file
		fileSystem[name] = file
	}
	file.Content = append(file.Content[:0], buffer...)
	file.ModTime = time.Now()
	return nil
}

func ReadFile(name string) ([]byte, error) {
	file, exists := fileSystem[name]
	if !exists || file.IsDir {
		return nil, errors.New("file does not exist or is a directory")
	}
	return file.Content, nil
}

func Open(name string) (*File, error) {
	file, exists := fileSystem[name]
	if !exists {
		return nil, errors.New("file does not exist")
	}
	fmt.Println("Opening file: " + name)
	return file, nil
}

func MkdirAll(path string, perm fs.FileMode) error {
	segments := strings.Split(path, "/")
	currentPath := "/"
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		nextPath := filepath.Join(currentPath, segment)
		if _, exists := fileSystem[nextPath]; !exists {
			dir := &File{Name: nextPath, IsDir: true, Perm: perm, ModTime: time.Now(), Entries: make(map[string]*File)}
			fileSystem[currentPath].Entries[segment] = dir
			fileSystem[nextPath] = dir
		}
		currentPath = nextPath
	}
	return nil
}

func Create(name string) (*File, error) {
	dirPath := filepath.Dir(name)
	dir, exists := fileSystem[dirPath]
	if !exists || !dir.IsDir {
		return nil, errors.New("directory does not exist " + name)
	}
	file := &File{Name: name, Perm: 0644, ModTime: time.Now()}
	dir.Entries[filepath.Base(name)] = file
	fileSystem[name] = file
	return file, nil
}

func IsNotExist(err error) bool {
	return errors.Is(err, errors.New("file does not exist"))
}

func RemoveAll(path string) error {
	delete(fileSystem, path)
	return nil
}

func Getwd() (string, error) {
	return "/", nil
}

func Remove(name string) error {
	delete(fileSystem, name)
	return nil
}

func Exit(code int) {
}

func Getenv(key string) string {
	return ""
}

type fileInfo struct {
	*File
}

func (fi fileInfo) Name() string       { return fi.File.Name }
func (fi fileInfo) Size() int64        { return int64(len(fi.File.Content)) }
func (fi fileInfo) Mode() fs.FileMode  { return fi.File.Perm }
func (fi fileInfo) ModTime() time.Time { return fi.File.ModTime }
func (fi fileInfo) IsDir() bool        { return fi.File.IsDir }
func (fi fileInfo) Sys() interface{}   { return nil }

type dirEntry struct {
	*File
}

func (de dirEntry) Name() string               { return de.File.Name }
func (de dirEntry) IsDir() bool                { return de.File.IsDir }
func (de dirEntry) Type() fs.FileMode          { return de.File.Perm }
func (de dirEntry) Info() (fs.FileInfo, error) { return fileInfo{de.File}, nil }
