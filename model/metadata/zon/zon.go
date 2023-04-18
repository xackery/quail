// zon is an EverQuest format that places zone terrain, objects, regions and lights
package zon

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
	"github.com/xackery/quail/pfs/archive"
)

// ZON is a zon file struct
type ZON struct {
	name          string
	version       uint32
	pfs           archive.ReadWriter
	models        []model
	objectManager *geo.ObjectManager
	regions       []region
	lights        []light

	terrains []*ter.TER
	mdses    []*mds.MDS
	mods     []*mod.MOD
}

type model struct {
	name     string
	baseName string
}

type region struct {
	name    string
	center  *geo.Vector3
	unknown *geo.Vector3
	extent  *geo.Vector3
}

type light struct {
	name     string
	position geo.Vector3
	color    geo.Vector3
	radius   float32
}

// New creates a new empty instance. Use NewFile to load an archive file on creation
func New(name string, pfs archive.ReadWriter) (*ZON, error) {
	if pfs == nil {
		return nil, fmt.Errorf("archive cannot be nil")
	}

	z := &ZON{
		name:          name,
		pfs:           pfs,
		objectManager: &geo.ObjectManager{},
	}
	return z, nil
}

// NewFile creates a new instance and loads provided file
func NewFile(name string, pfs archive.ReadWriter, file string) (*ZON, error) {
	e := &ZON{
		name: name,
		pfs:  pfs,
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

func (e *ZON) Name() string {
	return e.name
}

func (e *ZON) Data() []byte {
	w := bytes.NewBuffer(nil)
	err := e.Encode(w)
	if err != nil {
		log.Errorf("Failed to encode zon data: %s", err)
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

func (e *ZON) Regions() []region {
	return e.regions
}

func (e *ZON) Lights() []light {
	return e.lights
}

func (e *ZON) Models() []model {
	return e.models
}
