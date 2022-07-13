package ter

import "fmt"

func (e *TER) MaterialByID(id int) (string, error) {
	if id >= len(e.materials) {
		return "", fmt.Errorf("id '%d' is out of range (%d is max)", id, len(e.materials))
	}
	return e.materials[id].Name, nil
}
