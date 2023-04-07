package archive

import (
	"fmt"
)

type FileEntry struct {
	name string
	data []byte
}

func NewFileEntry(name string, data []byte) (*FileEntry, error) {
	fe := &FileEntry{}
	err := fe.SetName(name)
	if err != nil {
		return nil, fmt.Errorf("setname: %w", err)
	}
	err = fe.SetData(data)
	if err != nil {
		return nil, fmt.Errorf("setdata: %w", err)
	}
	return fe, nil
}

func (e *FileEntry) SetName(name string) error {
	e.name = name
	return nil
}

func (e *FileEntry) Name() string {
	return e.name
}

func (e *FileEntry) SetData(data []byte) error {
	e.data = data
	return nil
}

func (e *FileEntry) Data() []byte {
	return e.data
}
