//go:build tinygo.wasm || wasm

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

// Singletons
var jsFileSystem js.Value
var fileSystem = map[string]*File{
	"/": {Name: "/", IsDir: true, Entries: make(map[string]*File)},
}

// lastUpdateTimes tracks the last ModTime we pushed to JS for each file.
var lastUpdateTimes = make(map[string]time.Time)

func ExportFileSystem() js.Value {
	if jsFileSystem.IsUndefined() {
		jsFileSystem = js.Global().Get("Object").New()
		jsFileSystem.Set("write", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) != 2 {
				fmt.Println("Wanted 2 arguments for write")
				return false
			}
			src := args[0].String()
			buffer := make([]byte, args[1].Length())
			js.CopyBytesToGo(buffer, args[1])
			err := WriteFile(src, buffer, 0755)
			if err != nil {
				fmt.Println("Could not write file: " + err.Error())
				return false
			}
			return true
		}))
	}
	return jsFileSystem
}
func UpdateJSFile(fullPath string) {
	file, exists := fileSystem[fullPath]
	if !exists {
		removeJSFile(fullPath)
		return
	}

	lastTime, alreadyUpdated := lastUpdateTimes[fullPath]
	if alreadyUpdated && file.ModTime.Equal(lastTime) {
		return
	}
	lastUpdateTimes[fullPath] = file.ModTime
	dirPath := filepath.Dir(fullPath)
	dirObj := jsFileSystem.Get(dirPath)
	if dirObj.IsUndefined() {
		if dirPath == "/" {
			dirObj = jsFileSystem
		} else {
			dirObj = js.Global().Get("Object").New()
			dirPath = strings.TrimPrefix(dirPath, "/")
			jsFileSystem.Set(dirPath, dirObj)
		}
	}

	baseName := filepath.Base(fullPath)
	val := js.Global().Get("Uint8Array").New(len(file.Content))
	js.CopyBytesToJS(val, file.Content)
	dirObj.Set(baseName, val)

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
		return 0, fmt.Errorf("cannot write to a directory: %s", f.Name)
	}

	newSize := f.cursor + len(b)

	if newSize > len(f.Content) {
		growSize := len(f.Content) * 2
		if growSize < newSize {
			growSize = newSize
		}

		newContent := make([]byte, growSize)
		copy(newContent, f.Content)

		f.Content = nil

		f.Content = newContent
	}

	n := copy(f.Content[f.cursor:], b)
	f.cursor += n
	f.ModTime = time.Now()

	return n, nil
}

func (f *File) Close() error {
	UpdateJSFile(f.Name)
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
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}
	dir, exists := fileSystem[name]
	if !exists || !dir.IsDir {
		return nil, errors.New("read directory does not exist: " + name)
	}
	entries := make([]fs.DirEntry, 0, len(dir.Entries))
	for _, entry := range dir.Entries {
		entry.Name = filepath.Base(entry.Name)
		entries = append(entries, dirEntry{entry})
	}
	return entries, nil
}

func WriteFile(name string, buffer []byte, perm fs.FileMode) error {
	file, exists := fileSystem[name]
	if !exists {
		dirPath := filepath.Dir(name)
		dir, exists := fileSystem[dirPath]
		if exists && !dir.IsDir {
			return errors.New("directory is not a dir" + name + " dir path " + dirPath)
		}
		if !exists {
			MkdirAll(dirPath, 0755)
			dir = fileSystem[dirPath]
		}
		file = &File{Name: name, Perm: perm, ModTime: time.Now()}
		dir.Entries[filepath.Base(name)] = file
		fileSystem[name] = file
	}
	file.Content = append(file.Content[:0], buffer...)
	file.ModTime = time.Now()
	UpdateJSFile(name)
	return nil
}

func ReadFile(name string) ([]byte, error) {
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}
	file, exists := fileSystem[name]
	if !exists || file.IsDir {
		return nil, errors.New("file does not exist or is a directory" + name)
	}
	return file.Content, nil
}

func Open(name string) (*File, error) {
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}
	file, exists := fileSystem[name]
	if !exists {
		return nil, errors.New("file does not exist")
	}
	file.cursor = 0
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
	_, exists := fileSystem[name]
	if !exists {
		return nil
	}
	parentPath := filepath.Dir(name)
	parentFile, exists := fileSystem[parentPath]
	if exists && parentFile.IsDir {
		delete(parentFile.Entries, filepath.Base(name))
	}
	delete(fileSystem, name)

	// Remove from JS
	removeJSFile(name)
	return nil
}

func removeJSFile(fullPath string) {
	dirPath := filepath.Dir(fullPath)
	dirObj := jsFileSystem.Get(dirPath)
	if !dirObj.IsUndefined() {
		jsFileSystem.Delete(dirPath)
	}
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
