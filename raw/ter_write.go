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
		ter.NameAdd(material.EffectName)
		for _, prop := range material.Properties {
			ter.NameAdd(prop.Name)
			switch prop.Type {
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
		enc.Uint32(uint32(ter.NameIndex(material.EffectName)))
		enc.Uint32(uint32(len(material.Properties)))
		for _, prop := range material.Properties {
			enc.Uint32(uint32(ter.NameIndex(prop.Name)))
			enc.Uint32(uint32(prop.Type))
			switch prop.Type {
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
		enc.Float32(vertex.Position[0])
		enc.Float32(vertex.Position[1])
		enc.Float32(vertex.Position[2])
		enc.Float32(vertex.Normal[0])
		enc.Float32(vertex.Normal[1])
		enc.Float32(vertex.Normal[2])
		enc.Float32(vertex.Uv[0])
		enc.Float32(vertex.Uv[1])
		if ter.Version > 2 {
			enc.Float32(vertex.Uv2[0])
			enc.Float32(vertex.Uv2[1])
		}

	}

	for _, tri := range ter.Triangles {
		enc.Uint32(tri.Index[0])
		enc.Uint32(tri.Index[1])
		enc.Uint32(tri.Index[2])
		matID := int32(0)
		for _, mat := range ter.Materials {
			if mat.Name == tri.MaterialName {
				matID = mat.ID
				break
			}
		}
		enc.Int32(matID)
		enc.Uint32(tri.Flags)
	}

	err = enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
