package eqg

import "strings"

// Remove deletes an entry in an eqg, if any
func (e *EQG) Remove(name string) error {
	name = strings.ToLower(name)
	for i, f := range e.files {
		if strings.EqualFold(f.Name(), name) {
			e.files = append(e.files[:i], e.files[i+1:]...)
			return nil
		}
	}
	return nil
}
