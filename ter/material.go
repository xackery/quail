package ter

import "fmt"

func (e *TER) MaterialByID(id int32) (string, error) {
	if id == -1 {
		return "", nil
	}
	if int(id) >= len(e.materials) {
		return "", fmt.Errorf("id '%d' is out of range (%d is max)", id, len(e.materials))
	}
	return e.materials[id].Name, nil
}
