package ter

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

func (e *TER) Decode(r io.ReadSeeker) error {
	var err error

	modelName := strings.TrimSuffix(e.name, ".ter")

	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'E', 'Q', 'G', 'T'} {
		return fmt.Errorf("header does not match EQGM")
	}

	err = binary.Read(r, binary.LittleEndian, &e.version)
	if err != nil {
		return fmt.Errorf("read header version: %w", err)
	}
	dump.Hex(e.version, "version=%d", e.version)
	fmt.Println("version", e.version)

	switch e.version {
	case 2:
		err = e.loadVersion2(r, modelName)
		if err != nil {
			return fmt.Errorf("loadVersion2: %w", err)
		}
	case 3:
		err = e.loadVersion3(r, modelName)
		if err != nil {
			return fmt.Errorf("loadVersion3: %w", err)
		}
	default:
		return fmt.Errorf("unsupported *.zon version %d", e.version)
	}

	return nil
}

func (e *TER) loadVersion2(r io.Reader, modelName string) error {
	var err error
	nameLength := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &nameLength)
	if err != nil {
		return fmt.Errorf("read name length: %w", err)
	}
	dump.Hex(nameLength, "nameLength=%d", nameLength)

	materialCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &materialCount)
	if err != nil {
		return fmt.Errorf("read material count: %w", err)
	}
	dump.Hex(materialCount, "materialCount=%d", materialCount)

	verticesCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &verticesCount)
	if err != nil {
		return fmt.Errorf("read vertices count: %w", err)
	}
	dump.Hex(verticesCount, "verticesCount=%d", verticesCount)

	triangleCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &triangleCount)
	if err != nil {
		return fmt.Errorf("read triangle count: %w", err)
	}
	dump.Hex(triangleCount, "triangleCount=%d", triangleCount)

	/*err = binary.Read(r, binary.LittleEndian, uint32(len(e.boneAssignments)))
	if err != nil {
		return fmt.Errorf("read bone assignemt count: %w", err)
	}*/

	nameData := make([]byte, nameLength)

	err = binary.Read(r, binary.LittleEndian, &nameData)
	if err != nil {
		return fmt.Errorf("read nameData: %w", err)
	}

	names := make(map[uint32]string)

	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}
	dump.HexRange(nameData, int(nameLength), "nameData=(%d bytes, %d entries)", nameLength, len(names))
	for i := 0; i < int(materialCount); i++ {
		materialID := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read materialID: %w", err)
		}
		dump.Hex(materialID, "%dmaterialID=%d", i, materialID)

		nameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &nameOffset)
		if err != nil {
			return fmt.Errorf("read nameOffset: %w", err)
		}
		name, ok := names[nameOffset]
		if !ok {
			return fmt.Errorf("%dnames offset 0x%x not found", i, nameOffset)
		}
		dump.Hex(nameOffset, "%dnameOffset=0x%x(%s)", i, nameOffset, name)

		shaderOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &shaderOffset)
		if err != nil {
			return fmt.Errorf("read shaderOffset: %w", err)
		}
		shaderName, ok := names[shaderOffset]
		if !ok {
			return fmt.Errorf("%d names offset 0x%x not found", i, nameOffset)
		}
		dump.Hex(shaderOffset, "%dshaderOffset=0x%x(%s)", i, shaderOffset, shaderName)

		propertyCount := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &propertyCount)
		if err != nil {
			return fmt.Errorf("read propertyCount: %w", err)
		}
		dump.Hex(propertyCount, "%dpropertyCount=%d", i, propertyCount)

		err = e.MaterialManager.Add(name, shaderName)
		if err != nil {
			return fmt.Errorf("MaterialAdd %s: %w", name, err)
		}
		for j := 0; j < int(propertyCount); j++ {
			propertyNameOffset := uint32(0)
			err = binary.Read(r, binary.LittleEndian, &propertyNameOffset)
			if err != nil {
				return fmt.Errorf("read propertyNameOffset: %w", err)
			}
			propertyName, ok := names[propertyNameOffset]
			if !ok {
				return fmt.Errorf("%d%d read name offset: %w", i, j, err)
			}
			dump.Hex(propertyNameOffset, "%d%dpropertyNameOffset=0x%x(%s)", i, j, propertyNameOffset, propertyName)

			propertyType := uint32(0)
			err = binary.Read(r, binary.LittleEndian, &propertyType)
			if err != nil {
				return fmt.Errorf("read propertyType: %w", err)
			}
			dump.Hex(propertyType, "%d%dpropertyType=%d", i, j, propertyType)

			if propertyType == 0 {
				propFloatValue := float32(0)
				err = binary.Read(r, binary.LittleEndian, &propFloatValue)
				if err != nil {
					return fmt.Errorf("read propFloatValue: %w", err)
				}
				dump.Hex(propFloatValue, "%d%dpropertyFloat=%0.3f", i, j, propFloatValue)

				err = e.MaterialManager.PropertyAdd(name, propertyName, propertyType, fmt.Sprintf("%0.3f", propFloatValue))
				if err != nil {
					return fmt.Errorf("addMaterialProperty %s %s: %w", name, propertyName, err)
				}

			} else {
				propertyValue := uint32(0)
				err = binary.Read(r, binary.LittleEndian, &propertyValue)
				if err != nil {
					return fmt.Errorf("read propertyValue: %w", err)
				}
				dump.Hex(propertyValue, "%d%dpropertyValue=%d", i, j, propertyValue)
				propertyValueName := fmt.Sprintf("%d", propertyValue)
				if propertyType == 2 {
					propertyValueName, ok = names[propertyValue]
					if !ok {
						return fmt.Errorf("property %d names offset %d not found", j, propertyValue)
					}

					var data []byte
					data, err = e.archive.File(propertyValueName)
					if err != nil {
						data, err = e.archive.File(strings.ToLower(propertyValueName))
						if err != nil {
							fmt.Printf("warning: read material '%s' property '%s': %s\n", name, propertyName, err)
							//	return fmt.Errorf("read material via eqg %s: %w", propertyName, err)
						}
					}

					fe, err := archive.NewFileEntry(propertyValueName, data)
					if err != nil {
						return fmt.Errorf("new fileentry material %s: %w", propertyName, err)
					}

					e.files = append(e.files, fe)
				}
				err = e.MaterialManager.PropertyAdd(name, propertyName, propertyType, propertyValueName)
				if err != nil {
					return fmt.Errorf("addMaterialProperty %s %s: %w", name, propertyName, err)
				}
			}
		}
	}

	for i := 0; i < int(verticesCount); i++ {

		pos := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, pos)
		if err != nil {
			return fmt.Errorf("read vertex %d position: %w", i, err)
		}

		normal := &geo.Vector3{}
		err = binary.Read(r, binary.LittleEndian, normal)
		if err != nil {
			return fmt.Errorf("read vertex %d normal: %w", i, err)
		}

		tint := &geo.RGBA{R: 128, G: 128, B: 128, A: 1}

		uv := &geo.Vector2{}
		err = binary.Read(r, binary.LittleEndian, uv)
		if err != nil {
			return fmt.Errorf("read vertex %d uv: %w", i, err)
		}

		e.meshManager.VertexAdd(modelName, &geo.Vertex{
			Position: pos,
			Normal:   normal,
			Tint:     tint,
			Uv:       uv,
			Uv2:      uv,
		})

	}
	dump.HexRange([]byte{0x01, 0x02}, int(verticesCount)*32, "vertData=(%d bytes)", int(verticesCount)*32)

	for i := 0; i < int(triangleCount); i++ {
		pos := &geo.UIndex3{}
		err = binary.Read(r, binary.LittleEndian, pos)
		if err != nil {
			return fmt.Errorf("read triangle %d pos: %w", i, err)
		}

		materialID := int32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read triangle %d materialID: %w", i, err)
		}

		material, ok := e.MaterialManager.ByID(int(materialID))
		materialName := ""
		if ok {
			materialName = material.Name
		}

		flag := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &flag)
		if err != nil {
			return fmt.Errorf("read triangle %d flag: %w", i, err)
		}

		err = e.meshManager.TriangleAdd(modelName, pos, materialName, flag)
		if err != nil {
			return fmt.Errorf("triangleAdd %d: %w", i, err)
		}
	}
	dump.HexRange([]byte{0x03, 0x04}, int(triangleCount)*20, "triangleData=(%d bytes)", int(triangleCount)*20)
	return nil
}

func (e *TER) loadVersion3(r io.Reader, modelName string) error {
	var err error
	nameLength := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &nameLength)
	if err != nil {
		return fmt.Errorf("read name length: %w", err)
	}
	dump.Hex(nameLength, "nameLength=%d", nameLength)

	materialCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &materialCount)
	if err != nil {
		return fmt.Errorf("read material count: %w", err)
	}
	dump.Hex(materialCount, "materialCount=%d", materialCount)

	verticesCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &verticesCount)
	if err != nil {
		return fmt.Errorf("read vertices count: %w", err)
	}
	dump.Hex(verticesCount, "verticesCount=%d", verticesCount)

	triangleCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &triangleCount)
	if err != nil {
		return fmt.Errorf("read triangle count: %w", err)
	}
	dump.Hex(triangleCount, "triangleCount=%d", triangleCount)

	/*err = binary.Read(r, binary.LittleEndian, uint32(len(e.boneAssignments)))
	if err != nil {
		return fmt.Errorf("read bone assignemt count: %w", err)
	}*/

	nameData := make([]byte, nameLength)

	err = binary.Read(r, binary.LittleEndian, &nameData)
	if err != nil {
		return fmt.Errorf("read nameData: %w", err)
	}

	names := make(map[uint32]string)

	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[uint32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}
	dump.HexRange(nameData, int(nameLength), "nameData=(%d bytes, %d entries)", nameLength, len(names))
	for i := 0; i < int(materialCount); i++ {
		materialID := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read materialID: %w", err)
		}
		dump.Hex(materialID, "%dmaterialID=%d", i, materialID)

		nameOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &nameOffset)
		if err != nil {
			return fmt.Errorf("read nameOffset: %w", err)
		}
		name, ok := names[nameOffset]
		if !ok {
			return fmt.Errorf("%dnames offset 0x%x not found", i, nameOffset)
		}
		dump.Hex(nameOffset, "%dnameOffset=0x%x(%s)", i, nameOffset, name)

		shaderOffset := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &shaderOffset)
		if err != nil {
			return fmt.Errorf("read shaderOffset: %w", err)
		}
		shaderName, ok := names[shaderOffset]
		if !ok {
			return fmt.Errorf("%d names offset 0x%x not found", i, nameOffset)
		}
		dump.Hex(shaderOffset, "%dshaderOffset=0x%x(%s)", i, shaderOffset, shaderName)

		propertyCount := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &propertyCount)
		if err != nil {
			return fmt.Errorf("read propertyCount: %w", err)
		}
		dump.Hex(propertyCount, "%dpropertyCount=%d", i, propertyCount)

		err = e.MaterialManager.Add(name, shaderName)
		if err != nil {
			return fmt.Errorf("MaterialAdd %s: %w", name, err)
		}
		for j := 0; j < int(propertyCount); j++ {
			propertyNameOffset := uint32(0)
			err = binary.Read(r, binary.LittleEndian, &propertyNameOffset)
			if err != nil {
				return fmt.Errorf("read propertyNameOffset: %w", err)
			}
			propertyName, ok := names[propertyNameOffset]
			if !ok {
				return fmt.Errorf("%d%d read name offset: %w", i, j, err)
			}
			dump.Hex(propertyNameOffset, "%d%dpropertyNameOffset=0x%x(%s)", i, j, propertyNameOffset, propertyName)

			propertyType := uint32(0)
			err = binary.Read(r, binary.LittleEndian, &propertyType)
			if err != nil {
				return fmt.Errorf("read propertyType: %w", err)
			}
			dump.Hex(propertyType, "%d%dpropertyType=%d", i, j, propertyType)
			if propertyType == 0 {
				propFloatValue := float32(0)
				err = binary.Read(r, binary.LittleEndian, &propFloatValue)
				if err != nil {
					return fmt.Errorf("read propFloatValue: %w", err)
				}
				dump.Hex(propFloatValue, "%d%dpropertyFloat=%0.3f", i, j, propFloatValue)

				err = e.MaterialManager.PropertyAdd(name, propertyName, propertyType, fmt.Sprintf("%0.3f", propFloatValue))
				if err != nil {
					return fmt.Errorf("addMaterialProperty %s %s: %w", name, propertyName, err)
				}

			} else {
				propertyValue := uint32(0)
				err = binary.Read(r, binary.LittleEndian, &propertyValue)
				if err != nil {
					return fmt.Errorf("read propertyValue: %w", err)
				}
				dump.Hex(propertyValue, "%d%dpropertyValue=%d", i, j, propertyValue)
				propertyValueName := fmt.Sprintf("%d", propertyValue)
				if propertyType == 2 {
					propertyValueName, ok = names[propertyValue]
					if !ok {
						return fmt.Errorf("property %d names offset %d not found", j, propertyValue)
					}

					var data []byte
					if e.archive != nil {
						data, err = e.archive.File(propertyValueName)
						if err != nil {
							data, err = e.archive.File(strings.ToLower(propertyValueName))
							if err != nil {
								fmt.Printf("warning: read material '%s' property '%s': %s\n", name, propertyName, err)
								//	return fmt.Errorf("read material via eqg %s: %w", propertyName, err)
							}
						}
					}
					fe, err := archive.NewFileEntry(propertyValueName, data)
					if err != nil {
						return fmt.Errorf("new fileentry material %s: %w", propertyName, err)
					}

					e.files = append(e.files, fe)
				}
				err = e.MaterialManager.PropertyAdd(name, propertyName, propertyType, propertyValueName)
				if err != nil {
					return fmt.Errorf("addMaterialProperty %s %s: %w", name, propertyName, err)
				}
			}
		}
	}

	for i := 0; i < int(verticesCount); i++ {
		vert := geo.NewVertex()
		err = binary.Read(r, binary.LittleEndian, vert.Position)
		if err != nil {
			return fmt.Errorf("read vertex %d position: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, vert.Normal)
		if err != nil {
			return fmt.Errorf("read vertex %d normal: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, vert.Tint)
		if err != nil {
			return fmt.Errorf("read vertex %d tint: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, vert.Uv)
		if err != nil {
			return fmt.Errorf("read vertex %d uv: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, vert.Uv2)
		if err != nil {
			return fmt.Errorf("read vertex %d uv2: %w", i, err)
		}

		e.meshManager.VertexAdd(modelName, vert)
	}
	dump.HexRange([]byte{0x01, 0x02}, int(verticesCount)*32, "vertData=(%d bytes)", int(verticesCount)*32)

	for i := 0; i < int(triangleCount); i++ {
		pos := &geo.UIndex3{}
		err = binary.Read(r, binary.LittleEndian, pos)
		if err != nil {
			return fmt.Errorf("read triangle %d pos: %w", i, err)
		}

		materialID := int32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read triangle %d materialID: %w", i, err)
		}

		material, ok := e.MaterialManager.ByID(int(materialID))
		materialName := ""
		if ok {
			materialName = material.Name
		}

		flag := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &flag)
		if err != nil {
			return fmt.Errorf("read triangle %d flag: %w", i, err)
		}
		err = e.meshManager.TriangleAdd(modelName, pos, materialName, flag)
		if err != nil {
			return fmt.Errorf("addTriangle %d: %w", i, err)
		}
	}
	dump.HexRange([]byte{0x03, 0x04}, int(triangleCount)*20, "triangleData=(%d bytes)", int(triangleCount)*20)
	return nil
}