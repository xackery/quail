package def

// Zone is a zone
type Zone struct {
	Name    string
	Models  []string
	Objects []Object
	Regions []Region
	Lights  []Light
}

// Object is an object
type Object struct {
	Name      string
	ModelName string
	Position  Vector3
	Rotation  Vector3
	Scale     float32
}

// Region is a region
type Region struct {
	Name    string
	Center  Vector3
	Unknown Vector3
	Extent  Vector3
}

// Light is a light
type Light struct {
	Name     string
	Position Vector3
	Color    Vector3
	Radius   float32
}
