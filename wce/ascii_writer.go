package wce

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type AsciiWriteToken struct {
	basePath      string
	wce           *Wce
	lastWriter    io.WriteCloser
	lastWriterTag string
	writtenDefs   map[string]bool
	writers       map[string]io.WriteCloser
	writersUsed   map[string]bool
}

func NewAsciiWriteToken(path string, wce *Wce) *AsciiWriteToken {
	return &AsciiWriteToken{
		basePath:    path,
		wce:         wce,
		writtenDefs: make(map[string]bool),
		writers:     make(map[string]io.WriteCloser),
		writersUsed: make(map[string]bool),
	}
}

// TagIsWritten returns true if the tag was already written
func (a *AsciiWriteToken) TagIsWritten(tag string) bool {
	return a.writtenDefs[a.lastWriterTag+"-"+tag]
}

func (a *AsciiWriteToken) TagSetIsWritten(tag string) {
	a.writtenDefs[a.lastWriterTag+"-"+tag] = true
}

func (a *AsciiWriteToken) TagClearIsWritten() {
	for k := range a.writtenDefs {
		if strings.Contains(k, "_MDF_") {
			delete(a.writtenDefs, k)
		}
	}
}

func (a *AsciiWriteToken) Writer() (io.Writer, error) {
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

	isSubfolder := false
	w, ok := a.writers[tag]
	if !ok {
		rootFolder := tag

		if strings.Contains(rootFolder, "/") {
			chunks := strings.Split(rootFolder, "/")
			rootFolder = chunks[0]
			for i := 1; i < len(chunks)-1; i++ {
				rootFolder += "/" + chunks[i]
			}
			tag = chunks[len(chunks)-1]
			isSubfolder = true
		}

		path := filepath.Join(a.basePath, strings.ToLower(rootFolder), strings.ToLower(tag+".wce"))
		switch tag {
		case "world":
			path = filepath.Join(a.basePath, tag+".wce")
		}

		if isSubfolder {
			tag = rootFolder + "/" + tag
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
	a.lastWriterTag = tag
	a.lastWriter = w
	return nil
}

func (a *AsciiWriteToken) AddWriter(tag string, path string) error {
	_, ok := a.writers[tag]
	if ok {
		return fmt.Errorf("writer for tag %s already exists", tag)
	}

	dir := filepath.Dir(path)

	err := a.wce.FileSystem.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	w, err := a.wce.FileSystem.Create(path)
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
