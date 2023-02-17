// zon is an EverQuest format that places zone terrain, objects, regions and lights
package zon

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
)

// ZON is a zon file struct
type ZON struct {
	name    string
	archive common.ArchiveReadWriter
	models  []*model
	objects []*object
	regions []*region
	lights  []*light

	terrains []*ter.TER
	mdses    []*mds.MDS
	mods     []*mod.MOD
}

type model struct {
	name     string
	baseName string
}

type object struct {
	modelName   string
	name        string
	translation [3]float32
	rotation    [3]float32
	scale       float32
}

type region struct {
	name    string
	center  [3]float32
	unknown [3]float32
	extent  [3]float32
}

type light struct {
	name     string
	position [3]float32
	color    [3]float32
	radius   float32
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
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
	err = e.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return e, nil
}

func (e *ZON) Name() string {
	return e.name
}

func (e *ZON) Data() []byte {
	w := bytes.NewBuffer(nil)
	err := e.Encode(w)
	if err != nil {
		fmt.Println("failed to encode zon data:", err)
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

func (e *ZON) Regions() []*region {
	return e.regions
}

func (e *ZON) Lights() []*light {
	return e.lights
}

func (e *ZON) Objects() []*object {
	return e.objects
}

func (e *ZON) Models() []*model {
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
