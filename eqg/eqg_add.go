package eqg

import "fmt"

// Add adds a new entry to a eqg
func (e *EQG) Add(name string, data []byte) error {
	for _, f := range e.Files {
		if f.name == name {
			return fmt.Errorf("entry %s already exists", name)
		}
	}
	e.Files = append(e.Files, &fileEntry{name: name, data: data})
	return nil
}
