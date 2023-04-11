// ter is an EverQuest terrain model file
package ter

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/pfs/archive"
)

// TER is a terrain file struct
type TER struct {
	name            string
	materials       []*geo.Material
	version         uint32
	vertices        []*geo.Vertex
	triangles       []*geo.Triangle
	files           []archive.Filer
	archive         archive.ReadWriter
	particleRenders []*geo.ParticleRender
	particlePoints  []*geo.ParticlePoint
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.ReadWriter) (*TER, error) {
	t := &TER{
		name:    name,
		archive: pfs,
	}
	return t, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*TER, error) {
	e := &TER{
		name:    name,
		archive: pfs,
	}
	data, err := pfs.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", file, err)
	}
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return e, nil
}

func (e *TER) Name() string {
	return e.name
}

func (e *TER) Data() []byte {
	w := bytes.NewBuffer(nil)
	err := e.Encode(w)
	if err != nil {
		fmt.Println("failed to encode terrain data:", err)
		os.Exit(1)
	}
	return w.Bytes()
}

func (e *TER) SetName(value string) {
	e.name = value
}

func (e *TER) SetLayers(layers []*geo.Layer) error {
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

func (e *TER) SetParticleRenders(particles []*geo.ParticleRender) error {
	e.particleRenders = particles
	return nil
}

func (e *TER) SetParticlePoints(particles []*geo.ParticlePoint) error {
	e.particlePoints = particles
	return nil
}
