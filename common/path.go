package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Path struct {
	path      string
	files     []Filer
	fileCount int
}

// NewPath creates a new version of a path
func NewPath(path string) (*Path, error) {

	path = filepath.Dir(path)
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}
	e := &Path{path: path}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) < 1 {
			continue
		}
		if name[0] == '.' {
			continue
		}
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", path, name))
		if err != nil {
			return nil, err
		}
		fe, err := NewFileEntry(strings.ToLower(name), data)
		if err != nil {
			return nil, fmt.Errorf("newFileEntry '%s': %w", name, err)
		}
		e.files = append(e.files, fe)
	}
	e.fileCount = len(e.files)
	return e, nil
}

func (e *Path) File(name string) ([]byte, error) {
	for _, e := range e.files {
		if e.Name() == name {
			return e.Data(), nil
		}
	}
	return nil, fmt.Errorf("read %s: %w", name, os.ErrNotExist)
}

func (e *Path) Files() []Filer {
	return e.files
}

func (e *Path) Len() int {
	return e.fileCount
}

func (e *Path) String() string {
	return e.path
}

func (e *Path) WriteFile(name string, data []byte) error {
	w, err := os.Create(fmt.Sprintf("%s/%s", e.path, name))
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
