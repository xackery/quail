// zon is an EverQuest format that places zone terrain, objects, regions and lights
package zon

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xackery/quail/common"
)

// ZON is a zon file struct
type ZON struct {
	name    string
	path    string
	eqg     common.Archiver
	models  []*model
	objects []*Object
	regions []*Region
	lights  []*Light
}

type model struct {
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

func New(name string, archive common.Archiver) (*ZON, error) {
	z := &ZON{
		name: name,
		eqg:  archive,
	}
	return z, nil
}

func (e *ZON) SetPath(path string) {
	e.path = path
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

func (e *ZON) SetLayers(layers []*common.Layer) error {
	fmt.Println("TODO: set layers via zon")
	return nil
}

func (e *ZON) SetParticles(particles []*common.ParticleEntry) error {
	fmt.Println("TODO: set particles via zon")
	return nil
}
