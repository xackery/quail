package geo

// Skin is a skin
type Skin struct {
	Name                string
	InverseBindMatrices [][4]*Quad4
	Joints              map[int]*Joint
}

// Joint is a joint
type Joint struct {
	Name        string
	Children    []*Joint
	Translation *Vector3
}
