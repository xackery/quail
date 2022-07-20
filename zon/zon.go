// zon is an EverQuest format that places zone terrain, objects, regions and lights
package zon

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
)

// ZON is a zon file struct
type ZON struct {
	name    string
	archive common.ArchiveReadWriter
	models  []*Model
	objects []*Object
	regions []*Region
	lights  []*Light

	terrains []*ter.TER
	mdses    []*mds.MDS
	mods     []*mod.MOD
}

type Model struct {
	name     string
	baseName string
}

type Object struct {
	modelName   string
	name        string
	translation [3]float32
	rotation    [3]float32
	scale       float32
}

type Region struct {
	name    string
	center  [3]float32
	unknown [3]float32
	extent  [3]float32
}

type Light struct {
	name     string
	position [3]float32
	color    [3]float32
	radius   float32
}

func New(name string, archive common.ArchiveReadWriter) (*ZON, error) {
	if archive == nil {
		return nil, fmt.Errorf("archive cannot be nil")
	}

	z := &ZON{
		name:    name,
		archive: archive,
	}
	return z, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, archive common.ArchiveReadWriter, file string) (*ZON, error) {
	e := &ZON{
		name:    name,
		archive: archive,
	}
	data, err := archive.File(file)
	if err != nil {
		return nil, fmt.Errorf("file '%s': %w", file, err)
	}
	err = e.Load(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	return e, nil
}

func (e *ZON) Name() string {
	return e.name
}

func (e *ZON) Data() []byte {
	w := bytes.NewBuffer(nil)
	err := e.Save(w)
	if err != nil {
		fmt.Println("failed to save zon data:", err)
		os.Exit(1)
	}
	return w.Bytes()
}

// Models returns a slice of names
func (e *ZON) ModelNames() []string {
	names := []string{}
	for _, m := range e.models {
		names = append(names, m.name)
	}
	return names
}

func (e *ZON) Regions() []*Region {
	return e.regions
}

func (e *ZON) Lights() []*Light {
	return e.lights
}

func (e *ZON) Objects() []*Object {
	return e.objects
}

func (e *ZON) Models() []*Model {
	return e.models
}

func (e *ZON) SetLayers(layers []*common.Layer) error {
	fmt.Println("TODO: set layers via zon")
	return nil
}

func (e *ZON) SetParticleRenders(particles []*common.ParticleRender) error {
	fmt.Println("TODO: set particles via zon")
	return nil
}

func (e *ZON) SetParticlePoints(particles []*common.ParticlePoint) error {
	fmt.Println("TODO: set particles via zon")
	return nil
}
