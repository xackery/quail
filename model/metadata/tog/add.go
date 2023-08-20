package tog

import "github.com/xackery/quail/common"

func (e *TOG) AddObject(name string, position common.Vector3, rotation common.Vector3, scale float32, fileType string, fileName string) error {
	e.objects = append(e.objects, common.Object{
		Name:     name,
		Position: position,
		Rotation: rotation,
		Scale:    scale,
	})
	return nil
}
