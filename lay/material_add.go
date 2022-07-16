package lay

func (e *LAY) MaterialAdd(diffuseName string, normalName string) error {
	e.layers = append(e.layers, &layer{
		diffuse: diffuseName,
		normal:  normalName,
	})
	return nil
}
