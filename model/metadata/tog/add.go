package tog

import "github.com/xackery/quail/model/geo"

func (e *TOG) AddObject(name string, position *geo.Vector3, rotation *geo.Vector3, scale float32, fileType string, fileName string) error {
	e.objects = append(e.objects, &geo.Object{
		Name:     name,
		Position: position,
		Rotation: rotation,
		Scale:    scale,
		FileType: fileType,
		FileName: fileName,
	})
	return nil
}
