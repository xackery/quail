package mod

import "fmt"

func (e *MOD) MaterialByID(id int) (string, error) {
	if id >= len(e.materials) {
		return "", fmt.Errorf("id is out of range")
	}
	return e.materials[id].Name, nil
}
