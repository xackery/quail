package model

// Vector3 has X,Y,Z defined as float32
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// Vector2 has X,Y defined as float32
type Vector2 struct {
	X float32
	Y float32
}

// RGBA represents R,G,B,A as uint8
type RGBA struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// UIndex3 has X,Y,Z defined as uint32
type UIndex3 struct {
	X uint32
	Y uint32
	Z uint32
}

// Quad4  has X,Y,Z,W defined as float32
type Quad4 struct {
	X float32
	Y float32
	Z float32
	W float32
}
