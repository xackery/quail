package tog

import "github.com/g3n/engine/math32"

func (e *TOG) AddObject(name string, position math32.Vector3, rotation math32.Vector3, scale float32, fileType string, fileName string) error {
	e.objects = append(e.objects, &Object{
		name:     name,
		Position: position,
		Rotation: rotation,
		Scale:    scale,
		FileType: fileType,
		FileName: fileName,
	})
	return nil
}
