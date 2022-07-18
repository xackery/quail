package mds

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"sort"
	"strings"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/dump"
)

func (e *MDS) Load(r io.ReadSeeker) error {
	var err error
	header := [4]byte{}
	err = binary.Read(r, binary.LittleEndian, &header)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	dump.Hex(header, "header=%s", header)
	if header != [4]byte{'E', 'Q', 'G', 'S'} {
		return fmt.Errorf("header does not match EQGS")
	}

	version := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &version)
	if err != nil {
		return fmt.Errorf("read header version: %w", err)
	}
	dump.Hex(version, "version=%d", version)
	if version != 1 && version != 3 {
		return fmt.Errorf("version is %d, wanted 1 or 3", version)
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

	boneCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &boneCount)
	if err != nil {
		return fmt.Errorf("read bone count: %w", err)
	}
	dump.Hex(boneCount, "boneCount=%d", boneCount)

	subCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &subCount)
	if err != nil {
		return fmt.Errorf("read subCount: %w", err)
	}
	dump.Hex(subCount, "subCount=%d", subCount)

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
		name = strings.ToLower(name)
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

		if name == fmt.Sprintf("%s_02", e.name) {
			name = fmt.Sprintf("c_%s_s02_m01", e.name)
		}

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

				propertyValueName, ok := names[propertyValue]
				if !ok {
					return fmt.Errorf("property %d names offset %d not found", j, propertyValue)
				}

				data, err := e.archive.File(propertyValueName)
				if err != nil {
					fmt.Printf("warning: read material via eqg %s: %s\n", propertyName, err)
					//	return fmt.Errorf("read material via eqg %s: %w", propertyName, err)
				}
				fe, err := common.NewFileEntry(propertyValueName, data)
				if err != nil {
					return fmt.Errorf("new fileentry material %s: %w", propertyName, err)
				}
				e.files = append(e.files, fe)
				err = e.MaterialPropertyAdd(name, propertyName, propertyType, propertyValueName)
				if err != nil {
					return fmt.Errorf("addMaterialProperty %s %s: %w", name, propertyName, err)
				}
			}

		}
	}

	//64bytes worth
	for i := 0; i < int(boneCount); i++ {
		//52?
		materialID := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read bone %d materialID: %w", i, err)
		}
		dump.Hex(materialID, "%dmaterialid=%d(%s)", i, materialID, names[materialID])

		name := names[materialID]

		next := int32(0)
		err = binary.Read(r, binary.LittleEndian, &next)
		if err != nil {
			return fmt.Errorf("read bone %d next: %w", i, err)
		}
		dump.Hex(next, "%dnext=%d", i, next)

		childrenCount := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &childrenCount)
		if err != nil {
			return fmt.Errorf("read bone %d childrenCount: %w", i, err)
		}
		dump.Hex(childrenCount, "%dchildrenCount=%d", i, childrenCount)

		childIndex := int32(0)
		err = binary.Read(r, binary.LittleEndian, &childIndex)
		if err != nil {
			return fmt.Errorf("read bone %d childIndex: %w", i, err)
		}
		dump.Hex(childIndex, "%dchildIndex=%d", i, childIndex)

		pivot := math32.NewVec3()
		err = binary.Read(r, binary.LittleEndian, pivot)
		if err != nil {
			return fmt.Errorf("read bone %d pivot: %w", i, err)
		}
		dump.Hex(pivot, "%dpivot=%+v", i, pivot)

		rot := math32.NewVec4()
		err = binary.Read(r, binary.LittleEndian, rot)
		if err != nil {
			return fmt.Errorf("read bone %d rot: %w", i, err)
		}
		dump.Hex(rot, "%drot=%+v", i, rot)

		scale := math32.NewVec3()
		err = binary.Read(r, binary.LittleEndian, scale)
		if err != nil {
			return fmt.Errorf("read bone %d scale: %w", i, err)
		}
		dump.Hex(scale, "%dscale=%+v", i, scale)
		err = e.BoneAdd(name, next, childrenCount, childIndex, pivot, rot, scale)
		if err != nil {
			return fmt.Errorf("BoneAdd %d: %w", i, err)
		}
	}

	mainNameIndex := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &mainNameIndex)
	if err != nil {
		return fmt.Errorf("read mainNameIndex: %w", err)
	}
	dump.Hex(mainNameIndex, "mainNameIndex=%d", mainNameIndex)

	subNameIndex := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &subNameIndex)
	if err != nil {
		return fmt.Errorf("read subNameIndex: %w", err)
	}
	dump.Hex(subNameIndex, "subNameIndex=%d", subNameIndex)

	verticesCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &verticesCount)
	if err != nil {
		return fmt.Errorf("read vertices count: %w", err)
	}
	dump.Hex(verticesCount, "verticesCount=%d", verticesCount)

	faceCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &faceCount)
	if err != nil {
		return fmt.Errorf("read face count: %w", err)
	}
	dump.Hex(faceCount, "faceCount=%d", faceCount)

	boneAssignmentCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &boneAssignmentCount)
	if err != nil {
		return fmt.Errorf("read boneAssignmentCount: %w", err)
	}
	dump.Hex(boneAssignmentCount, "boneAssignmentCount=%d", boneAssignmentCount)

	for i := 0; i < int(verticesCount); i++ {

		pos := math32.NewVec3()
		err = binary.Read(r, binary.LittleEndian, pos)
		if err != nil {
			return fmt.Errorf("read vertex %d position: %w", i, err)
		}

		normal := math32.NewVec3()
		err = binary.Read(r, binary.LittleEndian, normal)
		if err != nil {
			return fmt.Errorf("read vertex %d normal: %w", i, err)
		}

		color := color.RGBA{}
		if version >= 3 {
			err = binary.Read(r, binary.LittleEndian, &color)
			if err != nil {
				return fmt.Errorf("read vertex %d color: %w", i, err)
			}

			unkUV := math32.Vector2{}
			err = binary.Read(r, binary.LittleEndian, &unkUV)
			if err != nil {
				return fmt.Errorf("read vertex %d unkUV: %w", i, err)
			}
		}

		uv := math32.NewVec2()
		err = binary.Read(r, binary.LittleEndian, uv)
		if err != nil {
			return fmt.Errorf("read vertex %d uv: %w", i, err)
		}
		uv.Y = -uv.Y
		tint := &common.Tint{R: 128, G: 128, B: 128}
		err = e.VertexAdd(pos, normal, tint, uv, uv)
		if err != nil {
			return fmt.Errorf("addVertex %d: %w", i, err)
		}
	}
	vSize := 32
	if version >= 3 {
		vSize += 12
	}
	dump.HexRange([]byte{0x01, 0x02}, int(verticesCount)*32, "vertData=(%d bytes)", int(verticesCount)*32)

	for i := 0; i < int(faceCount); i++ {
		pos := [3]uint32{}
		//pos := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("read face %d pos: %w", i, err)
		}

		materialID := int32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read face %d materialID: %w", i, err)
		}

		materialName, err := e.MaterialByID(int(materialID))
		if err != nil {
			return fmt.Errorf("material by id for face %d: %w", i, err)
		}

		flag := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &flag)
		if err != nil {
			return fmt.Errorf("read face %d flag: %w", i, err)
		}
		if materialName == "" {
			materialName = fmt.Sprintf("empty_%d", flag)
		}

		err = e.FaceAdd(pos, materialName, flag)
		if err != nil {
			return fmt.Errorf("faceAdd %d: %w", i, err)
		}
	}
	dump.HexRange([]byte{0x03, 0x04}, int(faceCount)*20, "faceData=(%d bytes)", int(faceCount)*20)
	sort.Sort(common.MaterialByName(e.materials))
	return nil
}
