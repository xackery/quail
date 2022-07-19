// ter is an EverQuest terrain model file
package ter

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/common"
)

// TER is a terrain file struct
type TER struct {
	name      string
	materials []*common.Material
	vertices  []*common.Vertex
	faces     []*common.Face
	files     []common.Filer
	eqg       common.Archiver
	particles []*common.ParticleEntry
}

func New(name string, archive common.Archiver) (*TER, error) {
	t := &TER{
		name: name,
		eqg:  archive,
	}
	return t, nil
}

func (e *TER) Name() string {
	return e.name
}

func (e *TER) Data() []byte {
	w := bytes.NewBuffer(nil)
	err := e.Save(w)
	if err != nil {
		fmt.Println("failed to save terrain data:", err)
		os.Exit(1)
	}
	return w.Bytes()
}

func (e *TER) SetName(value string) {
	e.name = value
}

func (e *TER) SetLayers(layers []*common.Layer) error {
	for _, o := range layers {
		err := e.MaterialAdd(o.Name, "")
		if err != nil {
			return fmt.Errorf("materialAdd: %w", err)
		}
		entry0Name := strings.ToLower(o.Entry0)
		entry1Name := strings.ToLower(o.Entry1)
		diffuseName := ""
		normalName := ""
		if strings.Contains(entry0Name, "_c.dds") {
			diffuseName = entry0Name
		}
		if strings.Contains(entry1Name, "_c.dds") {
			diffuseName = entry1Name
		}

		if strings.Contains(entry0Name, "_n.dds") {
			normalName = entry0Name
		}
		if strings.Contains(entry1Name, "_n.dds") {
			normalName = entry1Name
		}

		if len(diffuseName) > 0 {
			err = e.MaterialPropertyAdd(o.Name, "e_texturediffuse0", 2, diffuseName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", diffuseName, err)
			}
		}

		if len(normalName) > 0 {
			err = e.MaterialPropertyAdd(o.Name, "e_texturediffuse0", 2, normalName)
			if err != nil {
				return fmt.Errorf("materialPropertyAdd %s: %w", normalName, err)
			}
		}
	}
	return nil
}

func (e *TER) SetParticles(particles []*common.ParticleEntry) error {
	e.particles = particles
	return nil
}
