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
	a.writtenDefs = make(map[string]bool)
}

func (a *AsciiWriteToken) Writer() (*os.File, error) {
	if a.lastWriter == nil {
		return nil, fmt.Errorf("no writer set")
	}
	return a.lastWriter, nil
}

// Writer returns the file writer for a given tag
func (a *AsciiWriteToken) WriterByTag(tag string) (*os.File, error) {
	suffix := ""
	if strings.HasSuffix(tag, "_root") {
		suffix = "_root"
		tag = strings.TrimSuffix(tag, "_root")
	}
	if strings.HasSuffix(tag, "_ani") {
		suffix = "_ani"
		tag = strings.TrimSuffix(tag, "_ani")
	}

	rootTag := tag
	baseTag := tag
	if !a.wce.WorldDef.EqgVersion.Valid {
		rootTag = baseTagTrim(tag)
		baseTag = rootTag
	}
	baseTag = baseTag + suffix
	w, ok := a.writers[baseTag]
	if !ok {
		if len(rootTag) < 3 {
			return nil, fmt.Errorf("writer for short basetag %s (%s) does not exist", baseTag, tag)
		}
		rootTag = rootTag[:3]
		baseTag = rootTag + suffix

		w, ok = a.writers[baseTag]
		if !ok {
			w, ok = a.writers[a.wce.lastReadModelTag]
			if !ok {
				return nil, fmt.Errorf("writer for basetag %s (%s) does not exist (last read modeltag: %s)", baseTag, tag, a.wce.lastReadModelTag)
			}
		}
	}
	a.writersUsed[baseTag] = true
	return w, nil
}

func (a *AsciiWriteToken) SetWriter(tag string) error {
	w, err := a.WriterByTag(tag)
	if err != nil {
		return err
	}
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
