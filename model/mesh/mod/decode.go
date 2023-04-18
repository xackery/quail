package mod

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

// Decode decodes a MOD file
func (e *MOD) Decode(r io.ReadSeeker) error {
	var err error
	var ok bool

	modelName := strings.TrimSuffix(e.name, ".mod")

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQGM" {
		return fmt.Errorf("invalid header %s, wanted EQGM", header)
	}

	e.version = dec.Uint32()
	nameLength := int(dec.Uint32())
	materialCount := dec.Uint32()
	verticesCount := dec.Uint32()
	triangleCount := dec.Uint32()
	bonesCount := dec.Uint32()
	nameData := dec.Bytes(int(nameLength))

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

	log.Debugf("names: %+v", names)

	for i := 0; i < int(materialCount); i++ {
		material := geo.Material{}
		material.ID = dec.Int32()
		log.Debugf("material %d", material.ID)
		nameOffset := dec.Uint32()
		material.Name, ok = names[nameOffset]
		if !ok {
			return fmt.Errorf("material nameID %d not found", nameOffset)
		}
		shaderOffset := dec.Uint32()
		material.ShaderName, ok = names[shaderOffset]
		if !ok {
			return fmt.Errorf("material shader not found")
		}
		err = e.MaterialManager.Add(material)
		if err != nil {
			return fmt.Errorf("material add: %w", err)
		}

		propertyCount := dec.Uint32()
		for j := 0; j < int(propertyCount); j++ {
			property := geo.MaterialProperty{
				MaterialName: material.Name,
			}

			propertyNameOffset := dec.Uint32()
			property.Name, ok = names[propertyNameOffset]
			if !ok {
				return fmt.Errorf("material property name not found")
			}

			property.Category = dec.Uint32()
			if property.Category == 0 {
				property.Value = fmt.Sprintf("%0.3f", dec.Float32())
			} else {
				val := dec.Uint32()
				if property.Category == 2 {
					property.Value, ok = names[val]
					if !ok {
						return fmt.Errorf("material property value %d not found", val)
					}
				} else {
					property.Value = fmt.Sprintf("%d", val)
				}
			}

			err = e.MaterialManager.PropertyAdd(material.Name, property)
			if err != nil {
				return fmt.Errorf("material property add: %w", err)
			}
		}
	}

	for i := 0; i < int(verticesCount); i++ {
		v := geo.Vertex{}
		v.Position.X = dec.Float32()
		v.Position.Y = dec.Float32()
		v.Position.Z = dec.Float32()
		v.Normal.X = dec.Float32()
		v.Normal.Y = dec.Float32()
		v.Normal.Z = dec.Float32()
		if e.version <= 2 {
			v.Tint = geo.RGBA{R: 128, G: 128, B: 128, A: 255}
		} else {
			v.Tint = geo.RGBA{R: dec.Uint8(), G: dec.Uint8(), B: dec.Uint8(), A: dec.Uint8()}
		}
		v.Uv.X = dec.Float32()
		v.Uv.Y = dec.Float32()
		if e.version <= 2 {
			v.Uv2.X = 0
			v.Uv2.Y = 0
		} else {
			v.Uv2.X = dec.Float32()
			v.Uv2.Y = dec.Float32()
		}

		err = e.meshManager.VertexAdd(modelName, v)
		if err != nil {
			return fmt.Errorf("vertex add: %w", err)
		}
	}

	for i := 0; i < int(triangleCount); i++ {
		t := geo.Triangle{}
		t.Index.X = dec.Uint32()
		t.Index.Y = dec.Uint32()
		t.Index.Z = dec.Uint32()

		materialID := dec.Int32()
		material, ok := e.MaterialManager.ByID(materialID)
		if materialID != -1 && !ok {
			return fmt.Errorf("material %d not found", materialID)
		}
		t.MaterialName = material.Name

		t.Flag = dec.Uint32()
		e.meshManager.TriangleAdd(modelName, t)
	}

	for i := 0; i < int(bonesCount); i++ {
		bone := geo.Bone{}
		nameID := dec.Uint32()
		bone.Name, ok = names[nameID]
		if !ok {
			return fmt.Errorf("bone name %d not found", nameID)
		}
		bone.Next = dec.Int32()
		bone.ChildrenCount = dec.Uint32()
		bone.ChildIndex = dec.Int32()
		bone.Pivot.X = dec.Float32()
		bone.Pivot.Y = dec.Float32()
		bone.Pivot.Z = dec.Float32()
		bone.Rotation.X = dec.Float32()
		bone.Rotation.Y = dec.Float32()
		bone.Rotation.Z = dec.Float32()
		bone.Scale.X = dec.Float32()
		bone.Scale.Y = dec.Float32()
		bone.Scale.Z = dec.Float32()

		err = e.meshManager.BoneAdd(modelName, bone)
		if err != nil {
			return fmt.Errorf("bone add: %w", err)
		}
	}

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	log.Debugf("%s dencoded %d verts, %d triangles, %d bones, %d materials", e.name, e.meshManager.VertexTotalCount(), e.meshManager.TriangleTotalCount(), e.meshManager.BoneTotalCount(), e.MaterialManager.Count())

	return nil
}
