package eqg

// Remove deletes an entry in an eqg, if any
func (e *EQG) Remove(name string) error {
	for i, f := range e.files {
		if f.name == name {
			e.files = append(e.files[:i], e.files[i+1:]...)
			return nil
		}
	}
	return nil
}
