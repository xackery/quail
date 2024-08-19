package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/xackery/encdec"
)

// Model is a model
type Model struct {
	Header          *Header           `yaml:"header,omitempty"`
	FileType        string            `yaml:"file_type,omitempty"`
	Materials       []*Material       `yaml:"materials,omitempty"`
	ParticlePoints  []*ParticlePoint  `yaml:"particle_points,omitempty"`
	ParticleRenders []*ParticleRender `yaml:"particle_renders,omitempty"`
	Layers          []*Layer          `yaml:"layers,omitempty"`
	Bones           []Bone            `yaml:"bones,omitempty"`
	Vertices        []Vertex          `yaml:"vertices,omitempty"`
	Triangles       []Triangle        `yaml:"triangles,omitempty"`
}

// NewModel returns a new model
func NewModel(name string) *Model {
	return &Model{
		Header: &Header{
			Name: name,
		},
	}
}

func (v *Model) init() {
	if v.Header == nil {
		v.Header = &Header{}
	}
}

type Header struct {
	Name    string `yaml:"name"`
	Version int    `yaml:"version"`
}

// Material is a material
type Material struct {
	ID         int32               `yaml:"id,omitempty"`
	Name       string              `yaml:"name"`
	ShaderName string              `yaml:"shader_name"`
	Flag       uint32              `yaml:"flag"`
	Properties []*MaterialProperty `yaml:"properties,omitempty"`
	Animation  MaterialAnimation   `yaml:"animation,omitempty"`
}

// MaterialProperty is a material property
type MaterialProperty struct {
	Name     string `yaml:"name"`
	Category uint32 `yaml:"category"`
	Value    string `yaml:"value"`
	Data     []byte `yaml:"data,omitempty"`
}

func (m *MaterialProperty) String() string {
	return fmt.Sprintf("[%s, %d, %s]", m.Name, m.Category, m.Value)
}

type MaterialAnimation struct {
	Sleep    uint32   `yaml:"sleep"`
	Textures []string `yaml:"textures,omitempty"`
}

// Layer is a layer of materials to swap
type Layer struct {
	Material string `yaml:"material"`
	Diffuse  string `yaml:"diffuse"`
	Normal   string `yaml:"normal"`
}

func (v Layer) String() string {
	return fmt.Sprintf("[%s, %s, %s]", v.Material, v.Diffuse, v.Normal)
}

// Vector3 has X,Y,Z defined as float32
type Vector3 struct {
	X float32 `yaml:"fx"`
	Y float32 `yaml:"fy"`
	Z float32 `yaml:"fz"`
}

// String returns a string representation of the vector
func (v Vector3) String() string {
	return fmt.Sprintf("%f %f %f", v.X, v.Y, v.Z)
}

// Vector2 has X,Y defined as float32
type Vector2 struct {
	X float32 `yaml:"fx"`
	Y float32 `yaml:"fy"`
}

// Vertex is a vertex
type Vertex struct {
	Position Vector3 `yaml:"position"`
	Normal   Vector3 `yaml:"normal"`
	Tint     RGBA    `yaml:"tint"`
	Uv       Vector2 `yaml:"uv"`
	Uv2      Vector2 `yaml:"uv2"`
}

// RGBA represents R,G,B,A as uint8
type RGBA struct {
	R uint8 `yaml:"r"`
	G uint8 `yaml:"g"`
	B uint8 `yaml:"b"`
	A uint8 `yaml:"a"`
}

// Triangle is a triangle
type Triangle struct {
	Index        UIndex3 `yaml:"index"`
	MaterialName string  `yaml:"material_name"`
	Flag         uint32  `yaml:"flag"`
}

// UIndex3 has X,Y,Z defined as uint32
type UIndex3 struct {
	X uint32 `yaml:"ux"`
	Y uint32 `yaml:"uy"`
	Z uint32 `yaml:"uz"`
}

// Bone is a bone
type Bone struct {
	Name          string  `yaml:"name"`
	Next          int32   `yaml:"next"`
	ChildrenCount uint32  `yaml:"children_count"`
	ChildIndex    int32   `yaml:"child_index"`
	Pivot         Vector3 `yaml:"pivot"`
	Rotation      Quad4   `yaml:"rotation"`
	Scale         Vector3 `yaml:"scale"`
	Flags         uint32  `yaml:"flags"`
	Children      []int   `yaml:"children,omitempty"` // used by wld
}

// Quad4  has X,Y,Z,W defined as float32
type Quad4 struct {
	X float32 `yaml:"fx"`
	Y float32 `yaml:"fy"`
	Z float32 `yaml:"fz"`
	W float32 `yaml:"fw"`
}

// Normalize a quaternion
func Normalize(q Quad4) Quad4 {
	out := Quad4{}
	l := q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
	if l == 0 {
		out.X = 0
		out.Y = 0
		out.Z = 0
		out.W = 1
		return out
	}
	l = 1 / l
	out.X = q.X * l
	out.Y = q.Y * l
	out.Z = q.Z * l
	out.W = q.W * l
	return out
}

type BoneAnimation struct {
	Name       string                `yaml:"name"`
	FrameCount uint32                `yaml:"frame_count"`
	Frames     []*BoneAnimationFrame `yaml:"frames,omitempty"`
}

// BoneAnimationFrame is a bone animation frame
type BoneAnimationFrame struct {
	Milliseconds uint32  `yaml:"milliseconds"`
	Translation  Vector3 `yaml:"translation"`
	Rotation     Quad4   `yaml:"rotation"`
	Scale        Vector3 `yaml:"scale"`
}

// NameBuild prepares an EQG-styled name buffer list
func (model *Model) NameBuild(miscNames []string) (map[string]int32, []byte, error) {
	var err error

	names := make(map[string]int32)
	nameBuf := bytes.NewBuffer(nil)
	tmpNames := []string{}
	// append materials to tmpNames
	for _, o := range model.Materials {
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
	for _, bone := range model.Bones {
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

// VertexBuild prepares an EQG-styled vertex buffer list
func (model *Model) VertexBuild(version uint32, names map[string]int32) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)

	// verts
	for _, o := range model.Vertices {
		enc.Float32(o.Position.X)
		enc.Float32(o.Position.Y)
		enc.Float32(o.Position.Z)
		enc.Float32(o.Normal.X)
		enc.Float32(o.Normal.Y)
		enc.Float32(o.Normal.Z)
		if version <= 2 {
			enc.Float32(o.Uv.X)
			enc.Float32(o.Uv.Y)
		} else {
			enc.Uint8(o.Tint.R)
			enc.Uint8(o.Tint.G)
			enc.Uint8(o.Tint.B)
			enc.Uint8(o.Tint.A)
			enc.Float32(o.Uv.X)
			enc.Float32(o.Uv.Y)
			enc.Float32(o.Uv2.X)
			enc.Float32(o.Uv2.Y)
		}
	}
	return dataBuf.Bytes(), nil
}

// TriangleBuild prepares an EQG-styled triangle buffer list
func (model *Model) TriangleBuild(version uint32, names map[string]int32) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)

	// triangles
	for _, o := range model.Triangles {
		materialIdx := int32(-1)
		for idx, val := range model.Materials {
			if val.Name != o.MaterialName {
				continue
			}
			materialIdx = int32(idx)
			break
		}
		enc.Uint32(o.Index.X)
		enc.Uint32(o.Index.Y)
		enc.Uint32(o.Index.Z)
		enc.Int32(materialIdx)
		enc.Uint32(o.Flag)
	}

	return dataBuf.Bytes(), nil
}

// BoneBuild prepares an EQG-styled bone buffer list
func (model *Model) BoneBuild(version uint32, fileType string, names map[string]int32) ([]byte, error) {
	dataBuf := bytes.NewBuffer(nil)
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)

	// bones
	for _, o := range model.Bones {
		nameOffset := int32(-1)
		for key, val := range names {
			if key == o.Name {
				nameOffset = val
				break
			}
		}
		if nameOffset == -1 {
			return nil, fmt.Errorf("bone %s not found", o.Name)
		}

		enc.Int32(nameOffset)
		enc.Int32(o.Next)
		enc.Uint32(o.ChildrenCount)
		enc.Int32(o.ChildIndex)
		enc.Float32(o.Pivot.X)
		enc.Float32(o.Pivot.Y)
		enc.Float32(o.Pivot.Z)
		enc.Float32(o.Rotation.X)
		enc.Float32(o.Rotation.Y)
		enc.Float32(o.Rotation.Z)
		//enc.Float32(o.Rotation.W)
		enc.Float32(o.Scale.X)
		enc.Float32(o.Scale.Y)
		enc.Float32(o.Scale.Z)
		if fileType == "mod" {
			enc.Float32(1.0)
		}
	}
	return dataBuf.Bytes(), nil
}

// MaterialBuild prepares an EQG-styled material buffer list
func (model *Model) MaterialBuild(names map[string]int32) ([]byte, error) {
	var err error

	dataBuf := bytes.NewBuffer(nil)
	enc := encdec.NewEncoder(dataBuf, binary.LittleEndian)
	var nameOffset int32
	for materialID, o := range model.Materials {
		enc.Uint32(uint32(materialID))

		nameOffset = -1
		for key, offset := range names {
			if key == o.Name {
				nameOffset = offset
				break
			}
		}

		enc.Uint32(uint32(nameOffset))

		nameOffset = -1
		for key, offset := range names {
			if key == o.ShaderName {
				nameOffset = offset
				break
			}
		}
		if nameOffset == -1 {
			return nil, fmt.Errorf("shaderName %s not found", o.Name)
		}

		enc.Uint32(uint32(nameOffset))

		enc.Uint32(uint32(len(o.Properties)))

		for _, p := range o.Properties {
			nameOffset = -1
			for key, offset := range names {
				if key == p.Name {
					nameOffset = offset
					break
				}
			}
			if nameOffset == -1 {
				return nil, fmt.Errorf("%s prop %s not found", o.Name, p.Name)
			}

			enc.Uint32(uint32(nameOffset))
			enc.Uint32(p.Category)

			err = materialPropertyWrite(dataBuf, p.Value, names)
			if err != nil {
				return nil, fmt.Errorf("writePropertyValue: %w", err)
			}
		}
	}
	return dataBuf.Bytes(), nil
}

func materialPropertyWrite(buf *bytes.Buffer, value string, names map[string]int32) error {
	enc := encdec.NewEncoder(buf, binary.LittleEndian)
	val, err := strconv.Atoi(value)
	if err == nil {
		enc.Uint32(uint32(val))
		return nil
	}

	fVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		enc.Float32(float32(fVal))
		return nil
	}
	nameOffset := int32(-1)
	for key, offset := range names {
		if key == value {
			nameOffset = offset
			break
		}
	}
	if nameOffset == -1 {
		return fmt.Errorf("value %s: %w", value, err)
	}
	enc.Int32(nameOffset)
	return nil
}
