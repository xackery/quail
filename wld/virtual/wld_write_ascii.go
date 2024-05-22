package virtual

import (
	"fmt"
	"io"
)

func (wld *Wld) WriteAscii(w io.Writer) error {
	var err error

	for _, material := range wld.Materials {
		_, err = w.Write([]byte(material.Ascii()))
		if err != nil {
			return fmt.Errorf("write material %s: %w", material.Tag, err)
		}
	}

	for _, particleInstance := range wld.ParticleInstances {
		_, err = w.Write([]byte(particleInstance.Ascii()))
		if err != nil {
			return fmt.Errorf("write particleInstance %s: %w", particleInstance.Tag, err)
		}
	}

	return nil
}
