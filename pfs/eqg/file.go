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
	e.ContentsSummary = "\n"
	for i, fe := range e.files {
		base := float64(len(fe.Data()))
		out := ""
		num := float64(1024)
		if base < num*num*num*num {
			out = fmt.Sprintf("%0.0fG", base/num/num/num)
		}
		if base < num*num*num {
			out = fmt.Sprintf("%0.0fM", base/num/num)
		}
		if base < num*num {
			out = fmt.Sprintf("%0.0fK", base/num)
		}
		if base < num {
			out = fmt.Sprintf("%0.0fB", base)
		}
		e.ContentsSummary += fmt.Sprintf("%d %s:\t %s\n", i, out, fe.Name())
	}
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
