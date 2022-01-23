package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/g3n/engine/math32"
	"github.com/xackery/quail/common"
)

// Mesh information
type Mesh struct {
	hashIndex            uint32
	MaterialReference    uint32
	AnimationReference   uint32
	Center               math32.Vector3
	MaxDistance          float32
	MinPosition          math32.Vector3
	MaxPosition          math32.Vector3
	Verticies            []math32.Vector3
	TextureUVCoordinates []math32.Vector2
	Normals              []math32.Vector3
	Colors               []color.RGBA
	Indices              []*Polygon
}

type Polygon struct {
	IsSolid bool
	Vertex1 int
	Vertex2 int
	Vertex3 int
}

func LoadMesh(r io.ReadSeeker) (common.WldFragmenter, error) {
	v := &Mesh{}
	return v, nil
	err := parseMesh(r, v, false)
	if err != nil {
		return nil, fmt.Errorf("parse Mesh: %w", err)
	}
	return v, nil
}

func parseMesh(r io.ReadSeeker, v *Mesh, isNewWorldFormat bool) error {
	if v == nil {
		return fmt.Errorf("mesh is nil")
	}
	var value uint32
	err := binary.Read(r, binary.LittleEndian, &v.hashIndex)
	if err != nil {
		return fmt.Errorf("read hash index: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read flags: %w", err)
	}

	//if value != 0x00018003 && value != 0x00014003 {
	//	return fmt.Errorf("unknown mesh type, got 0x%x", value)
	//}//

	err = binary.Read(r, binary.LittleEndian, &v.MaterialReference)
	if err != nil {
		return fmt.Errorf("read material reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.AnimationReference)
	if err != nil {
		return fmt.Errorf("read animation reference: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown2: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center.X)
	if err != nil {
		return fmt.Errorf("read center x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center.Y)
	if err != nil {
		return fmt.Errorf("read center y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center.Z)
	if err != nil {
		return fmt.Errorf("read center z: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown1: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown2: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown3: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxDistance)
	if err != nil {
		return fmt.Errorf("read max distance: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center.X)
	if err != nil {
		return fmt.Errorf("read center x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center.Y)
	if err != nil {
		return fmt.Errorf("read center y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center.Z)
	if err != nil {
		return fmt.Errorf("read center z: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MinPosition.X)
	if err != nil {
		return fmt.Errorf("read min position x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MinPosition.Y)
	if err != nil {
		return fmt.Errorf("read min position y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MinPosition.Z)
	if err != nil {
		return fmt.Errorf("read min position z: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxPosition.X)
	if err != nil {
		return fmt.Errorf("read max position x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxPosition.Y)
	if err != nil {
		return fmt.Errorf("read max position y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxPosition.Z)
	if err != nil {
		return fmt.Errorf("read max position z: %w", err)
	}

	var vertexCount uint16
	err = binary.Read(r, binary.LittleEndian, &vertexCount)
	if err != nil {
		return fmt.Errorf("read vertex count: %w", err)
	}

	var textureCoordinateCount uint16
	err = binary.Read(r, binary.LittleEndian, &textureCoordinateCount)
	if err != nil {
		return fmt.Errorf("read texture coordinate count: %w", err)
	}

	var normalsCount uint16
	err = binary.Read(r, binary.LittleEndian, &normalsCount)
	if err != nil {
		return fmt.Errorf("read normals count: %w", err)
	}

	var colorsCount uint16
	err = binary.Read(r, binary.LittleEndian, &colorsCount)
	if err != nil {
		return fmt.Errorf("read colors count: %w", err)
	}

	var polygonCount uint16
	err = binary.Read(r, binary.LittleEndian, &polygonCount)
	if err != nil {
		return fmt.Errorf("read polygon count: %w", err)
	}

	var vertexPieceCount uint16
	err = binary.Read(r, binary.LittleEndian, &vertexPieceCount)
	if err != nil {
		return fmt.Errorf("read vertex piece count: %w", err)
	}

	var polygonTextureCount uint16
	err = binary.Read(r, binary.LittleEndian, &polygonTextureCount)
	if err != nil {
		return fmt.Errorf("read polygon texture count: %w", err)
	}

	var vertexTextureCount uint16
	err = binary.Read(r, binary.LittleEndian, &vertexTextureCount)
	if err != nil {
		return fmt.Errorf("read vertex texture count: %w", err)
	}

	var size9 uint16
	err = binary.Read(r, binary.LittleEndian, &size9)
	if err != nil {
		return fmt.Errorf("read size9: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return fmt.Errorf("read unknown2: %w", err)
	}

	scale := float32(1 / float32(int(1)<<value))

	for i := 0; i < int(vertexCount); i++ {
		pos := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &pos.X)
		if err != nil {
			return fmt.Errorf("read vertex x %d: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, &pos.Y)
		if err != nil {
			return fmt.Errorf("read vertex y %d: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, &pos.Z)
		if err != nil {
			return fmt.Errorf("read vertex z %d: %w", i, err)
		}
		pos.X *= scale
		pos.Y *= scale
		pos.Z *= scale

		v.Verticies = append(v.Verticies, pos)
	}

	for i := 0; i < int(textureCoordinateCount); i++ {
		if isNewWorldFormat {
			pos := math32.Vector2{}
			err = binary.Read(r, binary.LittleEndian, &pos.X)
			if err != nil {
				return fmt.Errorf("read texture coordinate 32 x %d: %w", i, err)
			}

			err = binary.Read(r, binary.LittleEndian, &pos.Y)
			if err != nil {
				return fmt.Errorf("read texture coordinate 32 y %d: %w", i, err)
			}
			pos.X /= 256
			pos.Y /= 256
			v.TextureUVCoordinates = append(v.TextureUVCoordinates, pos)

			continue
		}

		var tmpPos int16
		pos := math32.Vector2{}
		err = binary.Read(r, binary.LittleEndian, &tmpPos)
		if err != nil {
			return fmt.Errorf("read texture coordinate 16 x %d: %w", i, err)
		}

		pos.X = float32(tmpPos)

		err = binary.Read(r, binary.LittleEndian, &tmpPos)
		if err != nil {
			return fmt.Errorf("read texture coordinate 16 y %d: %w", i, err)
		}
		pos.Y = float32(tmpPos)
		pos.X /= 256
		pos.Y /= 256
		v.TextureUVCoordinates = append(v.TextureUVCoordinates, pos)
	}

	for i := 0; i < int(normalsCount); i++ {
		var val uint8
		pos := math32.Vector3{}
		err = binary.Read(r, binary.LittleEndian, &val)
		if err != nil {
			return fmt.Errorf("read normals x %d: %w", i, err)
		}
		pos.X = float32(val / 128)

		err = binary.Read(r, binary.LittleEndian, &val)
		if err != nil {
			return fmt.Errorf("read normals y %d: %w", i, err)
		}
		pos.Y = float32(val / 128)

		err = binary.Read(r, binary.LittleEndian, &val)
		if err != nil {
			return fmt.Errorf("read normals z %d: %w", i, err)
		}
		pos.Z = float32(val / 128)

		v.Normals = append(v.Normals, pos)
	}

	for i := 0; i < int(colorsCount); i++ {
		c := color.RGBA{}
		err = binary.Read(r, binary.LittleEndian, &c.R)
		if err != nil {
			return fmt.Errorf("read color r %d: %w", i, err)
		}
		err = binary.Read(r, binary.LittleEndian, &c.G)
		if err != nil {
			return fmt.Errorf("read color g %d: %w", i, err)
		}
		err = binary.Read(r, binary.LittleEndian, &c.B)
		if err != nil {
			return fmt.Errorf("read color b %d: %w", i, err)
		}
		err = binary.Read(r, binary.LittleEndian, &c.A)
		if err != nil {
			return fmt.Errorf("read color a %d: %w", i, err)
		}
		v.Colors = append(v.Colors, c)
	}

	for i := 0; i < int(polygonCount); i++ {
		var notSolidFlag int16
		err = binary.Read(r, binary.LittleEndian, &notSolidFlag)
		if err != nil {
			return fmt.Errorf("read notSolidFlag %d: %w", i, err)
		}
		p := &Polygon{}
		if notSolidFlag == 0 {
			//TODO: export separate collision flag
			p.IsSolid = true
		}
		err = binary.Read(r, binary.LittleEndian, &p.Vertex1)
		if err != nil {
			return fmt.Errorf("read vertex1 %d: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, &p.Vertex2)
		if err != nil {
			return fmt.Errorf("read vertex2 %d: %w", i, err)
		}

		err = binary.Read(r, binary.LittleEndian, &p.Vertex3)
		if err != nil {
			return fmt.Errorf("read vertex3 %d: %w", i, err)
		}

		v.Indices = append(v.Indices, p)

	}

	/*
	   MobPieces = new Dictionary<int, MobVertexPiece>();
	   int mobStart = 0;

	   for (int i = 0; i < vertexPieceCount; ++i)
	   {
	       int count = Reader.ReadInt16();
	       int index1 = Reader.ReadInt16();
	       var mobVertexPiece = new MobVertexPiece
	       {
	           Count = count,
	           Start = mobStart
	       };

	       mobStart += count;

	       MobPieces[index1] = mobVertexPiece;
	   }

	   MaterialGroups = new List<RenderGroup>();

	   StartTextureIndex = Int32.MaxValue;

	   for (int i = 0; i < polygonTextureCount; ++i)
	   {
	       var group = new RenderGroup();
	       group.PolygonCount = Reader.ReadUInt16();
	       group.MaterialIndex = Reader.ReadUInt16();
	       MaterialGroups.Add(group);

	       if (group.MaterialIndex < StartTextureIndex)
	       {
	           StartTextureIndex = group.MaterialIndex;
	       }
	   }

	   for (int i = 0; i < vertexTextureCount; ++i)
	   {
	       Reader.BaseStream.Position += 4;
	   }

	   for (int i = 0; i < size9; ++i)
	   {
	       Reader.BaseStream.Position += 12;
	   }

	   // In some rare cases, the number of uvs does not match the number of vertices
	   if (Vertices.Count != TextureUvCoordinates.Count)
	   {
	       int difference = Vertices.Count - TextureUvCoordinates.Count;

	       for (int i = 0; i < difference; ++i)
	       {
	           TextureUvCoordinates.Add(new vec2(0.0f, 0.0f));
	       }
	   }
	*/
	return nil
}

func (v *Mesh) FragmentType() string {
	return "Mesh"
}
func (e *Mesh) Data() []byte {
	buf := bytes.NewBuffer(nil)
	return buf.Bytes()
}
