// eqg is a pfs archive for EverQuest
package pfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/log"
)

// Pfs represents a modern everquest pfs archive
type Pfs struct {
	name            string
	files           []*FileEntry
	ContentsSummary string
	fileCount       int
}

// New creates a new empty instance. Use NewFile to load an archive on creation
func New(name string) (*Pfs, error) {
	e := &Pfs{
		name: name,
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	return e, nil
}

// NewFile takes path and loads it as an eqg archive
func NewFile(path string) (*Pfs, error) {
	e := &Pfs{
		name: filepath.Base(path),
	}
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	err = e.Read(r)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	return e, nil
}

// Remove deletes an entry in an eqg, if any
func (e *Pfs) Remove(name string) error {
	name = strings.ToLower(name)
	for i, f := range e.files {
		if !strings.EqualFold(f.Name(), name) {
			continue
		}

		e.files = append(e.files[:i], e.files[i+1:]...)
		return nil
	}
	return fmt.Errorf("file %s not found", name)
}

// Add adds a new entry to a eqg
func (e *Pfs) Add(name string, data []byte) error {
	name = strings.ToLower(name)

	if len(name) < 3 {
		return fmt.Errorf("name %s is too short", name)
	}

	for _, fe := range e.files {
		if fe.Name() == name {
			return nil
		}
	}

	log.Debugf("EQG adding %s (%d bytes)", name, len(data))
	e.files = append(e.files, NewFileEntry(name, data))
	return nil
}

// Extract places the pfs contents to path
func (e *Pfs) Extract(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("creating directory %s/\n", path)
			err = os.MkdirAll(path, 0766)
			if err != nil {
				return "", fmt.Errorf("mkdirall: %w", err)
			}
		}
		fi, err = os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("stat after mkdirall: %w", err)
		}
	}
	if !fi.IsDir() {
		return "", fmt.Errorf("%s is not a directory", path)
	}

	extractStdout := ""
	for i, file := range e.files {
		err = os.WriteFile(fmt.Sprintf("%s/%s", path, file.Name()), file.Data(), 0644)
		if err != nil {
			return "", fmt.Errorf("index %d: %w", i, err)
		}
		extractStdout += file.Name() + ", "
	}
	if len(e.files) == 0 {
		return "", fmt.Errorf("no files found to extract")
	}
	extractStdout = extractStdout[0 : len(extractStdout)-2]
	return fmt.Sprintf("extracted %d file%s to %s: %s", len(e.files), helper.Pluralize(len(e.files)), path, extractStdout), nil
}

// File returns data of a file
func (e *Pfs) File(name string) ([]byte, error) {
	for _, f := range e.files {
		if f.Name() == name || strings.EqualFold(f.Name(), strings.ToLower(name)) {
			return f.Data(), nil
		}
	}
	return nil, fmt.Errorf("read %s: %w", name, os.ErrNotExist)
}

func (e *Pfs) Close() error {
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

func (e *Pfs) Len() int {
	return len(e.files)
}

// Files returns a string array of every file inside an EQG
func (e *Pfs) Files() []*FileEntry {
	return e.files
}

func (e *Pfs) SetFile(name string, data []byte) error {
	name = strings.ToLower(name)
	if len(name) < 3 {
		return fmt.Errorf("name %s is too short", name)
	}
	for _, file := range e.files {
		if strings.EqualFold(file.Name(), name) {
			return file.SetData(data)
		}
	}
	e.files = append(e.files, NewFileEntry(name, data))
	return nil
}
