package raw

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/xackery/encdec"
)

// Write writes a ter file
func (ter *Ter) Write(w io.Writer) error {
	var err error
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQGT")
	enc.Uint32(ter.Version)

	ter.NameClear()

	for _, material := range ter.Materials {
		ter.NameAdd(material.Name)
		ter.NameAdd(material.ShaderName)
		for _, prop := range material.Properties {
			ter.NameAdd(prop.Name)
			switch prop.Category {
			case 2:
				ter.NameAdd(prop.Value)
			default:
			}
		}
	}

	nameData := ter.NameData()
	enc.Uint32(uint32(len(nameData))) // nameLength
	enc.Uint32(uint32(len(ter.Materials)))
	enc.Uint32(uint32(len(ter.Vertices)))
	enc.Uint32(uint32(len(ter.Triangles)))
	enc.Bytes(nameData)

	for _, material := range ter.Materials {
		enc.Int32(material.ID)
		enc.Uint32(uint32(ter.NameIndex(material.Name)))
		enc.Uint32(uint32(ter.NameIndex(material.ShaderName)))
		enc.Uint32(uint32(len(material.Properties)))
		for _, prop := range material.Properties {
			enc.Uint32(uint32(ter.NameIndex(prop.Name)))
			enc.Uint32(uint32(prop.Category))
			switch prop.Category {
			case 0:
				fval, err := strconv.ParseFloat(prop.Value, 32)
				if err != nil {
					return err
				}
				enc.Float32(float32(fval))
			case 2:
				enc.Int32(ter.NameIndex(prop.Value))
			default:
				return err
			}
		}
	}

	for _, vertex := range ter.Vertices {
		enc.Float32(vertex.Position.X)
		enc.Float32(vertex.Position.Y)
		enc.Float32(vertex.Position.Z)
		enc.Float32(vertex.Normal.X)
		enc.Float32(vertex.Normal.Y)
		enc.Float32(vertex.Normal.Z)
		enc.Float32(vertex.Uv.X)
		enc.Float32(vertex.Uv.Y)
		if ter.Version > 2 {
			enc.Float32(vertex.Uv2.X)
			enc.Float32(vertex.Uv2.Y)
		}

	}

	for _, tri := range ter.Triangles {
		enc.Uint32(tri.Index.X)
		enc.Uint32(tri.Index.Y)
		enc.Uint32(tri.Index.Z)
		matID := int32(0)
		for _, mat := range ter.Materials {
			if mat.Name == tri.MaterialName {
				matID = mat.ID
				break
			}
		}
		enc.Int32(matID)
		enc.Uint32(tri.Flag)
	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
