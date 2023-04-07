package geo

// Object is a generic object used in .TOG files
type Object struct {
	Name     string
	Position *Vector3
	Rotation *Vector3
	Scale    float32
	FileType string
	FileName string
}
