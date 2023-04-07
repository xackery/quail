package lay

import "github.com/xackery/quail/model/geo"

func (e *LAY) MaterialAdd(name string, diffuseName string, normalName string) error {
	e.layers = append(e.layers, &geo.Layer{
		Name:   name,
		Entry0: diffuseName,
		Entry1: normalName,
	})
	return nil
}
