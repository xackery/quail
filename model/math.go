package model

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

// RGBA represents R,G,B,A as uint8
type RGBA struct {
	R uint8 `yaml:"r"`
	G uint8 `yaml:"g"`
	B uint8 `yaml:"b"`
	A uint8 `yaml:"a"`
}

// UIndex3 has X,Y,Z defined as uint32
type UIndex3 struct {
	X uint32 `yaml:"ux"`
	Y uint32 `yaml:"uy"`
	Z uint32 `yaml:"uz"`
}

// Quad4  has X,Y,Z,W defined as float32
type Quad4 struct {
	X float32 `yaml:"fx"`
	Y float32 `yaml:"fy"`
	Z float32 `yaml:"fz"`
	W float32 `yaml:"fw"`
}
