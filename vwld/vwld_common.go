package vwld

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
