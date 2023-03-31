package mod

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
)

func (e *MOD) Decode(r io.ReadSeeker) error {
	var err error
	e.isSkinned = false

	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'E', 'Q', 'G', 'M'} {
		return fmt.Errorf("header does not match EQGM")
	}

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("read header version: %w", err)
	}
	dump.Hex(version, "version=%d", version)
	if version > 3 {
		return fmt.Errorf("version is %d, wanted < 4", version)
	}

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

	boneCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &boneCount)
	if err != nil {
		return fmt.Errorf("read bone count: %w", err)
	}
	dump.Hex(boneCount, "boneCount=%d", boneCount)

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

		err = e.MaterialAdd(name, shaderName)
		if err != nil {
			return fmt.Errorf("addMaterial %s: %w", name, err)
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
				dump.Hex(propFloatValue, "%d%dpropertyFloat=%0.2f", i, j, propFloatValue)

				err = e.MaterialPropertyAdd(name, propertyName, propertyType, fmt.Sprintf("%0.2f", propFloatValue))
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
						return fmt.Errorf("%d%d material %s property offset %d not found", i, j, name, propertyNameOffset)
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
					} else {
						data, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", e.path, propertyValueName))
						if err != nil {
							fmt.Printf("warning: read material via %s: %s\n", propertyName, err.Error())
							//return fmt.Errorf("read material via path %s: %w", propertyName, err)
						}
					}
					fe, err := common.NewFileEntry(propertyValueName, data)
					if err != nil {
						return fmt.Errorf("new fileentry material %s: %w", propertyName, err)
					}

					e.files = append(e.files, fe)
				}
				err = e.MaterialPropertyAdd(name, propertyName, propertyType, propertyValueName)
				if err != nil {
					return fmt.Errorf("addMaterialProperty %s %s: %w", name, propertyName, err)
				}
			}

		}
	}

	for i := 0; i < int(verticesCount); i++ {

		vertex := &common.Vertex{}

		err = binary.Read(r, binary.LittleEndian, &vertex.Position)
		if err != nil {
			return fmt.Errorf("read vertex %d position: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, &vertex.Normal)
		if err != nil {
			return fmt.Errorf("read vertex %d normal: %w", i, err)
		}

		if version < 3 {

			uv := [2]float32{}
			err = binary.Read(r, binary.LittleEndian, &uv)
			if err != nil {
				return fmt.Errorf("read vertex %d uv: %w", i, err)
			}

			vertex.Tint = [4]uint8{128, 128, 128, 0}
		} else {
			// TODO: may be misaligned (RGB vs RGBA)
			err = binary.Read(r, binary.LittleEndian, &vertex.Tint)
			if err != nil {
				return fmt.Errorf("read vertex %d tint: %w", i, err)
			}

			err = binary.Read(r, binary.LittleEndian, &vertex.Uv)
			if err != nil {
				return fmt.Errorf("read vertex %d uv: %w", i, err)
			}

			err = binary.Read(r, binary.LittleEndian, &vertex.Uv2)
			if err != nil {
				return fmt.Errorf("read vertex %d uv2: %w", i, err)
			}
		}

		e.vertices = append(e.vertices, vertex)
	}
	vSize := 32
	if version >= 3 {
		vSize += 12
	}
	dump.HexRange([]byte{0x01, 0x02}, int(verticesCount)*32, "vertData=(%d bytes)", int(verticesCount)*32)

	for i := 0; i < int(triangleCount); i++ {
		pos := [3]uint32{}
		//pos := [3]float32{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("read triangle %d pos: %w", i, err)
		}

		materialID := int32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read triangle %d materialID: %w", i, err)
		}

		materialName, err := e.MaterialByID(int(materialID))
		if err != nil {
			return fmt.Errorf("material by id for triangle %d: %w", i, err)
		}

		flag := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &flag)
		if err != nil {
			return fmt.Errorf("read triangle %d flag: %w", i, err)
		}

		if materialName == "" {
			materialName = fmt.Sprintf("empty_%d", flag)
		}
		err = e.TriangleAdd(pos, materialName, flag)
		if err != nil {
			return fmt.Errorf("triangleAdd %d: %w", i, err)
		}
	}
	dump.HexRange([]byte{0x03, 0x04}, int(triangleCount)*20, "triangleData=(%d bytes)", int(triangleCount)*20)

	//64bytes worth
	for i := 0; i < int(boneCount); i++ {
		//52?

		//	vertex := e.vertices[i]

		bone := &common.Bone{}

		materialID := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read bone %d materialID: %w", i, err)
		}
		dump.Hex(materialID, "%dmaterialid=%d(%s)", i, materialID, names[materialID])

		bone.Name = names[materialID]

		err = binary.Read(r, binary.LittleEndian, &bone.Next)
		if err != nil {
			return fmt.Errorf("read bone %d next: %w", i, err)
		}
		dump.Hex(bone.Next, "%dnext=%d", i, bone.Next)

		err = binary.Read(r, binary.LittleEndian, &bone.ChildrenCount)
		if err != nil {
			return fmt.Errorf("read bone %d childrenCount: %w", i, err)
		}
		dump.Hex(bone.ChildrenCount, "%dchildrenCount=%d", i, bone.ChildrenCount)

		err = binary.Read(r, binary.LittleEndian, &bone.ChildIndex)
		if err != nil {
			return fmt.Errorf("read bone %d childIndex: %w", i, err)
		}
		dump.Hex(bone.ChildIndex, "%dchildIndex=%d", i, bone.ChildIndex)

		err = binary.Read(r, binary.LittleEndian, &bone.Pivot)
		if err != nil {
			return fmt.Errorf("read bone %d pivot: %w", i, err)
		}
		dump.Hex(bone.Pivot, "%dpivot=%+v", i, bone.Pivot)

		err = binary.Read(r, binary.LittleEndian, &bone.Rotation)
		if err != nil {
			return fmt.Errorf("read bone %d rot: %w", i, err)
		}
		dump.Hex(bone.Rotation, "%drot=%+v", i, bone.Rotation)

		err = binary.Read(r, binary.LittleEndian, &bone.Scale)
		if err != nil {
			return fmt.Errorf("read bone %d scale: %w", i, err)
		}
		dump.Hex(bone.Scale, "%dscale=%+v", i, bone.Scale)

		e.bones = append(e.bones, bone)
	}
	return nil
}
