package def

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// Mesh is a mesh
type Mesh struct {
	Name            string
	FileType        string
	Vertices        []Vertex
	Triangles       []Triangle
	Bones           []Bone
	Materials       []*Material
	ParticlePoints  []*ParticlePoint
	ParticleRenders []*ParticleRender
}

// Material is a material
type Material struct {
	ID         int32
	Name       string
	ShaderName string
	Flag       uint32
	Properties []*MaterialProperty
}

// MaterialProperty is a material property
type MaterialProperty struct {
	Name     string
	Category uint32
	Value    string
	Data     []byte
}

// Vector3 has X,Y,Z defined as float32
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// String returns a string representation of the vector
func (v Vector3) String() string {
	return fmt.Sprintf("%f %f %f", v.X, v.Y, v.Z)
}

// Vector2 has X,Y defined as float32
type Vector2 struct {
	X float32
	Y float32
}

// Vertex is a vertex
type Vertex struct {
	Position Vector3
	Normal   Vector3
	Tint     RGBA
	Uv       Vector2
	Uv2      Vector2
}

// RGBA represents R,G,B,A as uint8
type RGBA struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// Triangle is a triangle
type Triangle struct {
	Index        UIndex3
	MaterialName string
	Flag         uint32
}

// UIndex3 has X,Y,Z defined as uint32
type UIndex3 struct {
	X uint32
	Y uint32
	Z uint32
}

// Bone is a bone
type Bone struct {
	Name          string
	Next          int32
	ChildrenCount uint32
	ChildIndex    int32
	Pivot         Vector3
	Rotation      Quad4
	Scale         Vector3
}

// Quad4  has X,Y,Z,W defined as float32
type Quad4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

type BoneAnimation struct {
	Name       string
	FrameCount uint32
	Frames     []*BoneAnimationFrame
}

// BoneAnimationFrame is a bone animation frame
type BoneAnimationFrame struct {
	Milliseconds uint32
	Translation  *Vector3
	Rotation     *Quad4
	Scale        *Vector3
}

// NameBuild prepares an EQG-styled name buffer list
func (mesh *Mesh) nameBuild(miscNames []string) (map[string]int32, []byte, error) {
	var err error

	names := make(map[string]int32)
	nameBuf := bytes.NewBuffer(nil)
	tmpNames := []string{}
	// append materials to tmpNames
	for _, o := range mesh.Materials {
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
	for _, bone := range mesh.Bones {
		tmpNames = append(tmpNames, bone.Name)
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

// ApplyQuaternion transforms this vector by multiplying it by
// the specified quaternion and then by the quaternion inverse.
// It basically applies the rotation encoded in the quaternion to this vector.
func ApplyQuaternion(v Vector3, q Quad4) Vector3 {
	x := v.X
	y := v.Y
	z := v.Z

	qx := q.X
	qy := q.Y
	qz := q.Z
	qw := q.W

	// calculate quat * vector
	ix := qw*x + qy*z - qz*y
	iy := qw*y + qz*x - qx*z
	iz := qw*z + qx*y - qy*x
	iw := -qx*x - qy*y - qz*z
	// calculate result * inverse quat
	v.X = ix*qw + iw*-qx + iy*-qz - iz*-qy
	v.Y = iy*qw + iw*-qy + iz*-qx - ix*-qz
	v.Z = iz*qw + iw*-qz + ix*-qy - iy*-qx
	return v
}

func (mesh *Mesh) Encode(version uint32, w io.Writer) error {
	switch mesh.FileType {
	case "mod":
		return mesh.MODEncode(version, w)
	case "mds":
		return mesh.MDSEncode(version, w)
	case "ter":
		return mesh.TEREncode(version, w)
	}
	return fmt.Errorf("unknown file type: %s", mesh.FileType)
}
