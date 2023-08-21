package common

import (
	"bytes"
	"fmt"
)

// Zone is a zone
type Zone struct {
	Version int
	Name    string
	Models  []string
	Objects []Object
	Regions []Region
	Lights  []Light
	Lits    []*RGBA
}

// Object is an object
type Object struct {
	Name      string
	ModelName string
	Position  Vector3
	Rotation  Vector3
	Scale     float32
}

// Region is a region
type Region struct {
	Name    string
	Center  Vector3
	Unknown Vector3
	Extent  Vector3
}

// Light is a light
type Light struct {
	Name     string
	Position Vector3
	Color    Vector3
	Radius   float32
}

// NameBuild prepares an EQG-styled name buffer list
func (zone *Zone) NameBuild(miscNames []string) (map[string]int32, []byte, error) {
	var err error

	names := make(map[string]int32)
	nameBuf := bytes.NewBuffer(nil)
	tmpNames := []string{}

	for _, name := range miscNames {
		isNew := true
		for key := range names {
			if key == name {
				isNew = false
				break
			}
		}
		if !isNew {
			continue
		}

		tmpNames = append(tmpNames, name)
	}

	// append materials to tmpNames
	for _, o := range zone.Objects {
		tmpNames = append(tmpNames, o.Name)
	}

	// append regions to tmpNames
	for _, r := range zone.Regions {
		tmpNames = append(tmpNames, r.Name)
	}

	// append lights to tmpNames
	for _, l := range zone.Lights {
		tmpNames = append(tmpNames, l.Name)
	}

	for _, name := range tmpNames {
		isNew := true
		for key := range names {
			if key == name {
				isNew = false
				break
			}
		}
		if !isNew {
			continue
		}

		names[name] = int32(nameBuf.Len())

		_, err = nameBuf.Write([]byte(name))
		if err != nil {
			return nil, nil, fmt.Errorf("write name: %w", err)
		}
		_, err = nameBuf.Write([]byte{0})
		if err != nil {
			return nil, nil, fmt.Errorf("write 0: %w", err)
		}
	}

	return names, nameBuf.Bytes(), nil
}
