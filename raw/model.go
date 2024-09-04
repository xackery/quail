package raw

// Bone is a bone
type Bone struct {
	Name          string
	Next          int32
	ChildrenCount uint32
	ChildIndex    int32
	Pivot         [3]float32
	Rotation      [4]float32
	Scale         [3]float32
	Scale2        float32
}

// Triangle is a triangle
type Triangle struct {
	Index        [3]uint32
	MaterialName string
	Flag         uint32
}

type Material struct {
	ID         int32
	Name       string
	ShaderName string
	Flag       uint32
	Properties []*MaterialProperty
	Animation  MaterialAnimation
}

// MaterialProperty is a material property
type MaterialProperty struct {
	Name     string
	Category uint32
	Value    string
	Data     []byte
}

type MaterialAnimation struct {
	Sleep    uint32
	Textures []string
}

// Vertex is a vertex
type Vertex struct {
	Position [3]float32
	Normal   [3]float32
	Tint     [4]uint8
	Uv       [2]float32
	Uv2      [2]float32
}
