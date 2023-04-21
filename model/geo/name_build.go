package geo

import (
	"bytes"
	"fmt"
	"strconv"
)

// NameBuild prepares an EQG-styled name buffer list
func NameBuild(matManager *MaterialManager, meshManager *MeshManager, miscNames []string) (map[string]int32, []byte, error) {
	var err error

	names := make(map[string]int32)
	nameBuf := bytes.NewBuffer(nil)
	tmpNames := []string{}
	// append materials to tmpNames
	for _, o := range matManager.materials {
		tmpNames = append(tmpNames, o.Name)
		tmpNames = append(tmpNames, o.ShaderName)
		for _, p := range o.Properties {
			tmpNames = append(tmpNames, p.Name)
			_, err = strconv.Atoi(p.Value)
			if err != nil {
				_, err = strconv.ParseFloat(p.Value, 64)
				if err != nil {
					tmpNames = append(tmpNames, p.Value)
				}
			}
		}
	}

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

	// append bones to tmpNames
	for _, mesh := range meshManager.Meshes() {
		for _, bone := range mesh.Bones {
			tmpNames = append(tmpNames, bone.Name)
		}
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
