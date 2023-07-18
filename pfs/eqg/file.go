package eqg

import (
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/pfs/archive"
)

// File returns data of a file
func (e *EQG) File(name string) ([]byte, error) {
	for _, f := range e.files {
		if f.Name() == name || strings.EqualFold(f.Name(), strings.ToLower(name)) {
			return f.Data(), nil
		}
	}
	return nil, fmt.Errorf("read %s: %w", name, os.ErrNotExist)
}

func (e *EQG) Close() error {
	e.files = nil
	e.name = ""
	e.fileCount = 0
	return nil
}

func (e *EQG) Len() int {
	return len(e.files)
}

// Files returns a string array of every file inside an EQG
func (e *EQG) Files() []archive.Filer {
	return e.files
}

func (e *EQG) WriteFile(name string, data []byte) error {
	//name = strings.ToLower(name)
	for _, file := range e.files {
		if file.Name() == name {
			return file.SetData(data)
		}
	}
	fe, err := archive.NewFileEntry(name, data)
	if err != nil {
		return fmt.Errorf("newFileEntry: %w", err)
	}
	e.files = append(e.files, fe)
	return nil
}
