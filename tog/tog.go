package tog

import "github.com/g3n/engine/math32"

type TOG struct {
	name    string
	objects []*Object
}

type Object struct {
	Name     string
	Position math32.Vector3
	Rotation math32.Vector3
	Scale    float32
	FileType string
	FileName string
}

func New(name string) (*TOG, error) {
	e := &TOG{
		name: name,
	}
	return e, nil
}
