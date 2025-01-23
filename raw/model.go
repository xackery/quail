package raw

// Bone is a bone
type Bone struct {
	Name          string
	Next          int32
	ChildrenCount uint32
	ChildIndex    int32
	Pivot         [3]float32
	Quaternion    [4]float32
	Scale         [3]float32
}

// Face is a triangle
type Face struct {
	Index        [3]uint32
	MaterialName string
	Flags        uint32
}

type Material struct {
	ID         int32
	Name       string
	EffectName string
	Flag       uint32
	Properties []*MaterialParam
	Animation  MaterialAnimation
}

type MaterialParamType uint32

const (
	MaterialParamTypeUnused MaterialParamType = iota
	MaterialParamTypeInt
	MaterialParamTypeTexture
	MaterialParamTypeColor
)

// MaterialParam is a material property
type MaterialParam struct {
	Name  string
	Type  MaterialParamType
	Value string
	Data  []byte
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
