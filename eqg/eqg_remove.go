package eqg

// Remove deletes an entry in an eqg, if any
func (e *EQG) Remove(name string) error {
	for i, f := range e.Files {
		if f.name == name {
			e.Files = append(e.Files[:i], e.Files[i+1:]...)
			return nil
		}
	}
	return nil
}
