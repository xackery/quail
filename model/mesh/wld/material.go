package wld

import "fmt"

// MaterialByID returns a material by id
func (e *WLD) MaterialByID(id int) (string, error) {
	if id == -1 {
		return "", nil
	}
	if id >= len(e.materials) {
		return "", fmt.Errorf("id '%d' is out of range (%d is max)", id, len(e.materials))
	}
	return e.materials[id].Name, nil
}
