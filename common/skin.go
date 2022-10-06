package common

type Skin struct {
	Name                string
	InverseBindMatrices [][4][4]float32
	Joints              map[int]*Joint
}

type Joint struct {
	Name        string
	Children    []*Joint
	Translation [3]float32
}
