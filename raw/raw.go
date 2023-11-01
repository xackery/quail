package raw

import (
	"fmt"
	"io"
)

type Reader interface {
	Read(r io.ReadSeeker) error
	SetFileName(name string)
}

type Writer interface {
	FileName() string
	Write(w io.Writer) error
}

type ReadWriter interface {
	Reader
	Writer
}

// Vector3 has X,Y,Z defined as float32
type Vector3 struct {
	X float32 `yaml:"fx"`
	Y float32 `yaml:"fy"`
	Z float32 `yaml:"fz"`
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
	Scale2        float32 `yaml:"scale2"`
}

// Quad4  has X,Y,Z,W defined as float32
type Quad4 struct {
	X float32 `yaml:"fx"`
	Y float32 `yaml:"fy"`
	Z float32 `yaml:"fz"`
	W float32 `yaml:"fw"`
}

type Material struct {
	ID         int32               `yaml:"id"`
	Name       string              `yaml:"name"`
	ShaderName string              `yaml:"shader_name"`
	Flag       uint32              `yaml:"flag"`
	Properties []*MaterialProperty `yaml:"properties"`
	Animation  MaterialAnimation   `yaml:"animation"`
}

// MaterialProperty is a material property
type MaterialProperty struct {
	Name     string `yaml:"name"`
	Category uint32 `yaml:"category"`
	Value    string `yaml:"value"`
	Data     []byte `yaml:"data,omitempty"`
}

type MaterialAnimation struct {
	Sleep    uint32   `yaml:"sleep"`
	Textures []string `yaml:"textures,omitempty"`
}

// New takes an extension and returns a ReadWriter that can parse it
func New(ext string) ReadWriter {
	switch ext {
	case ".ani":
		return &Ani{}
	case ".bmp":
		return &Bmp{}
	case ".dat":
		return &Dat{}
	case ".dds":
		return &Dds{}
	case ".edd":
		return &Edd{}
	case ".lay":
		return &Lay{}
	case ".lit":
		return &Lit{}
	case ".lod":
		return &Lod{}
	case ".mds":
		return &Mds{}
	case ".mod":
		return &Mod{}
	case ".png":
		return &Png{}
	case ".prt":
		return &Prt{}
	case ".pts":
		return &Pts{}
	case ".ter":
		return &Ter{}
	case ".tog":
		return &Tog{}
	case ".wld":
		return &Wld{}
	case ".zon":
		return &Zon{}
	default:
		return nil
	}
}

// Read takes an extension and a reader and returns a ReadWriter that can parse it
func Read(ext string, r io.ReadSeeker) (ReadWriter, error) {
	reader := New(ext)
	if reader == nil {
		return nil, fmt.Errorf("unknown extension %s", ext)
	}
	err := reader.Read(r)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

// Write takes an extension and a writer and returns a ReadWriter that can parse it
func Write(ext string, w io.Writer) (ReadWriter, error) {
	writer := New(ext)
	if writer == nil {
		return nil, fmt.Errorf("unknown extension %s", ext)
	}
	err := writer.Write(w)
	if err != nil {
		return nil, err
	}
	return writer, nil
}
