package tog

func (e *TOG) AddObject(name string, position [3]float32, rotation [3]float32, scale float32, fileType string, fileName string) error {
	e.objects = append(e.objects, &Object{
		Name:     name,
		Position: position,
		Rotation: rotation,
		Scale:    scale,
		FileType: fileType,
		FileName: fileName,
	})
	return nil
}
