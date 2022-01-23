package mod

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/dump"
)

func (e *MOD) Load(r io.ReadSeeker) error {
	var err error

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

		err = e.AddMaterial(name, shaderName)
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

			propertyValue := uint32(0)
			err = binary.Read(r, binary.LittleEndian, &propertyValue)
			if err != nil {
				return fmt.Errorf("read propertyValue: %w", err)
			}
			dump.Hex(propertyValue, "%d%dpropertyValue=%d", i, j, propertyValue)
			err = e.AddMaterialProperty(name, propertyName, propertyType, float32(propertyValue), propertyValue)
			if err != nil {
				return fmt.Errorf("addMaterialProperty %s %s: %w", name, propertyName, err)
			}
		}
	}

	for i := 0; i < int(verticesCount); i++ {

		pos := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("read vertex %d position: %w", i, err)
		}

		normal := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &normal)
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

		uv := math32.Vector2{}
		err = binary.Read(r, binary.LittleEndian, &uv)
		if err != nil {
			return fmt.Errorf("read vertex %d uv: %w", i, err)
		}
		err = e.AddVertex(pos, normal, uv)
		if err != nil {
			return fmt.Errorf("addVertex %d: %w", i, err)
		}
	}
	vSize := 32
	if version >= 3 {
		vSize += 12
	}
	dump.HexRange([]byte{0x01, 0x02}, int(verticesCount)*32, "vertData=(%d bytes)", int(verticesCount)*32)

	for i := 0; i < int(triangleCount); i++ {
		pos := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("read triangle %d pos: %w", i, err)
		}

		materialID := uint32(0)
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
		err = e.AddTriangle(pos, materialName, flag)
		if err != nil {
			return fmt.Errorf("addTriangle %d: %w", i, err)
		}
	}
	dump.HexRange([]byte{0x03, 0x04}, int(triangleCount)*20, "triangleData=(%d bytes)", int(triangleCount)*20)

	//64bytes worth
	for i := 0; i < int(boneCount); i++ {
		//52?
		materialID := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &materialID)
		if err != nil {
			return fmt.Errorf("read bone %d materialID: %w", i, err)
		}

		dump.Hex(materialID, "%dmaterialid=%d(%s)", i, materialID, names[materialID])
		maxVal := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &maxVal)
		if err != nil {
			return fmt.Errorf("read bone %d maxVal: %w", i, err)
		}
		dump.Hex(maxVal, "%dmaxVal=%d", i, maxVal)

		fiveVal := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &fiveVal)
		if err != nil {
			return fmt.Errorf("read bone %d fiveVal: %w", i, err)
		}
		dump.Hex(fiveVal, "%dfiveVal=%d", i, fiveVal)

		oneVal := uint32(0)
		err = binary.Read(r, binary.LittleEndian, &oneVal)
		if err != nil {
			return fmt.Errorf("read bone %d oneVal: %w", i, err)
		}
		dump.Hex(oneVal, "%doneVal=%d", i, oneVal)

		pos := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("read bone %d pos: %w", i, err)
		}
		dump.Hex(pos, "%dpos=%+v", i, pos)

		rot := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &rot)
		if err != nil {
			return fmt.Errorf("read bone %d rot: %w", i, err)
		}
		dump.Hex(rot, "%drot=%+v", i, rot)

		scale := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &scale)
		if err != nil {
			return fmt.Errorf("read bone %d scale: %w", i, err)
		}
		dump.Hex(scale, "%dscale=%+v", i, scale)

		chunk := make([]byte, 60)
		err = binary.Read(r, binary.LittleEndian, &chunk)
		if err != nil {
			return fmt.Errorf("read chunk %d: %w", i, err)
		}
		dump.Hex(chunk, "%dchunk=(%d bytes)", i, len(chunk))
		// pure f's

		//5
		//1
	}
	return nil
}
