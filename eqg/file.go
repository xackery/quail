package eqg

import "fmt"

// File returns data of a file
func (e *EQG) File(name string) ([]byte, error) {
	for _, f := range e.files {
		if f.Name() == name {
			return f.Data(), nil
		}
	}
	return nil, fmt.Errorf("%s not found", name)
}
