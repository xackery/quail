package fragment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"

	"github.com/xackery/quail/common"
)

// Mesh information
type Mesh struct {
	name                 string
	MaterialReference    uint32
	AnimationReference   uint32
	Center               [3]float32
	MaxDistance          float32
	MinPosition          [3]float32
	MaxPosition          [3]float32
	Verticies            [][3]float32
	TextureUVCoordinates [][2]float32
	Normals              [][3]float32
	Colors               []color.RGBA
	Indices              []*Polygon
}

type Polygon struct {
	IsSolid bool
	Vertex1 int16
	Vertex2 int16
	Vertex3 int16
}

func LoadMesh(r io.ReadSeeker) (common.WldFragmenter, error) {
	v := &Mesh{}
	err := parseMesh(r, v, false)
	if err != nil {
		return nil, fmt.Errorf("parse Mesh: %w", err)
	}
	return v, nil
}

func parseMesh(r io.ReadSeeker, v *Mesh, isNewWorldFormat bool) error {
	var err error
	if v == nil {
		return fmt.Errorf("mesh is nil")
	}
	var value uint32
	v.name, err = nameFromHashIndex(r)
	if err != nil {
		return fmt.Errorf("nameFromHasIndex: %w", err)
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

	err = binary.Read(r, binary.LittleEndian, v.Center)
	if err != nil {
		return fmt.Errorf("read center: %w", err)
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

	err = binary.Read(r, binary.LittleEndian, &v.Center[0])
	if err != nil {
		return fmt.Errorf("read center x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center[1])
	if err != nil {
		return fmt.Errorf("read center y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.Center[2])
	if err != nil {
		return fmt.Errorf("read center z: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MinPosition[0])
	if err != nil {
		return fmt.Errorf("read min position x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MinPosition[1])
	if err != nil {
		return fmt.Errorf("read min position y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MinPosition[2])
	if err != nil {
		return fmt.Errorf("read min position z: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxPosition[0])
	if err != nil {
		return fmt.Errorf("read max position x: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxPosition[1])
	if err != nil {
		return fmt.Errorf("read max position y: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &v.MaxPosition[2])
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
		pos := [3]float32{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("read vertex %d: %w", i, err)
		}

		pos[0] *= scale
		pos[1] *= scale
		pos[2] *= scale

		v.Verticies = append(v.Verticies, pos)
	}

	for i := 0; i < int(textureCoordinateCount); i++ {
		if isNewWorldFormat {
			pos := [2]float32{}
			err = binary.Read(r, binary.LittleEndian, &pos)
			if err != nil {
				return fmt.Errorf("read texture coordinate 32 %d: %w", i, err)
			}

			pos[0] /= 256
			pos[1] /= 256
			v.TextureUVCoordinates = append(v.TextureUVCoordinates, pos)

			continue
		}

		pos := [2]float32{}
		err = binary.Read(r, binary.LittleEndian, &pos)
		if err != nil {
			return fmt.Errorf("read texture coordinate 16 %d: %w", i, err)
		}

		pos[0] /= 256
		pos[1] /= 256
		v.TextureUVCoordinates = append(v.TextureUVCoordinates, pos)
	}

	for i := 0; i < int(normalsCount); i++ {
		var val uint8
		pos := [3]float32{}
		err = binary.Read(r, binary.LittleEndian, &val)
		if err != nil {
			return fmt.Errorf("read normals x %d: %w", i, err)
		}
		pos[0] = float32(val / 128)

		err = binary.Read(r, binary.LittleEndian, &val)
		if err != nil {
			return fmt.Errorf("read normals y %d: %w", i, err)
		}
		pos[1] = float32(val / 128)

		err = binary.Read(r, binary.LittleEndian, &val)
		if err != nil {
			return fmt.Errorf("read normals z %d: %w", i, err)
		}
		pos[2] = float32(val / 128)

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
