package lay

func (e *LAY) MaterialAdd(name string, diffuseName string, normalName string) error {
	e.layers = append(e.layers, &layer{
		name:    name,
		diffuse: diffuseName,
		normal:  normalName,
	})
	return nil
}
