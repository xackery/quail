package raw

import "github.com/xackery/quail/model"

// Bone is a bone
type Bone struct {
	Name          string        `yaml:"name"`
	Next          int32         `yaml:"next"`
	ChildrenCount uint32        `yaml:"children_count"`
	ChildIndex    int32         `yaml:"child_index"`
	Pivot         model.Vector3 `yaml:"pivot"`
	Rotation      model.Quad4   `yaml:"rotation"`
	Scale         model.Vector3 `yaml:"scale"`
	Scale2        float32       `yaml:"scale2"`
}

// Triangle is a triangle
type Triangle struct {
	Index        model.UIndex3 `yaml:"index"`
	MaterialName string        `yaml:"material_name"`
	Flag         uint32        `yaml:"flag"`
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

// Vertex is a vertex
type Vertex struct {
	Position model.Vector3 `yaml:"position"`
	Normal   model.Vector3 `yaml:"normal"`
	Tint     model.RGBA    `yaml:"tint"`
	Uv       model.Vector2 `yaml:"uv"`
	Uv2      model.Vector2 `yaml:"uv2"`
}
