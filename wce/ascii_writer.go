package wce

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AsciiWriteToken struct {
	basePath    string
	wce         *Wce
	lastWriter  *os.File
	writtenDefs map[string]bool
	writers     map[string]*os.File
	writersUsed map[string]bool
}

func NewAsciiWriteToken(path string, wce *Wce) *AsciiWriteToken {
	return &AsciiWriteToken{
		basePath:    path,
		wce:         wce,
		writtenDefs: make(map[string]bool),
		writers:     make(map[string]*os.File),
		writersUsed: make(map[string]bool),
	}
}

// TagIsWritten returns true if the tag was already written
func (a *AsciiWriteToken) TagIsWritten(tag string) bool {
	return a.writtenDefs[tag]
}

func (a *AsciiWriteToken) TagSetIsWritten(tag string) {
	a.writtenDefs[tag] = true
}

func (a *AsciiWriteToken) TagClearIsWritten() {
	for k := range a.writtenDefs {
		if strings.Contains(k, "_MDF_") {
			delete(a.writtenDefs, k)
		}
	}
}

func (a *AsciiWriteToken) Writer() (*os.File, error) {
	if a.lastWriter == nil {
		return nil, fmt.Errorf("no writer set")
	}
	return a.lastWriter, nil
}

func (a *AsciiWriteToken) SetWriter(tag string) error {
	var err error
	if tag == "" {
		tag = "world"
	}

	w, ok := a.writers[tag]
	if !ok {
		rootFolder := tag
		if strings.Contains(rootFolder, "_") {
			rootFolder = strings.Split(rootFolder, "_")[0]
		}
		path := filepath.Join(a.basePath, strings.ToLower(rootFolder), strings.ToLower(tag+".wce"))
		switch tag {
		case "world", "region":
			path = filepath.Join(a.basePath, tag+".wce")
		}

		err = a.AddWriter(tag, path)
		if err != nil {
			return err
		}
		w, ok = a.writers[tag]
		if !ok {
			return fmt.Errorf("writer for tag %s not found", tag)
		}
	}

	a.writersUsed[tag] = true

	a.lastWriter = w
	return nil
}

func (a *AsciiWriteToken) AddWriter(tag string, path string) error {
	_, ok := a.writers[tag]
	if ok {
		return fmt.Errorf("writer for tag %s already exists", tag)
	}

	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	w, err := os.Create(path)
	if err != nil {
		return err
	}

	a.wce.writeAsciiHeader(w)
	a.writers[tag] = w
	return nil
}

func (a *AsciiWriteToken) IsWriterUsed(tag string) bool {
	return a.writersUsed[tag]
}

func (a *AsciiWriteToken) Close() {
	for _, writer := range a.writers {
		writer.Close()
	}
}
