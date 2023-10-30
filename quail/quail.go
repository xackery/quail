package quail

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
	"github.com/xackery/quail/model/metadata/ani"
	"github.com/xackery/quail/model/metadata/lay"
	"github.com/xackery/quail/model/metadata/lit"
	"github.com/xackery/quail/model/metadata/prt"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/zon"
)

type Quail struct {
	Header                 *common.Header
	Models                 []*common.Model
	Animations             []*common.Animation
	Zone                   *common.Zone
	materialCache          map[string]*common.Material
	IsExtensionVersionDump bool
}

// New returns a new Quail instance
func New() *Quail {
	return &Quail{
		materialCache: make(map[string]*common.Material),
	}
}

// Close flushes any memory and closes any open files
func (e *Quail) Close() error {
	e.Models = nil
	e.Animations = nil
	e.Zone = nil
	e.materialCache = make(map[string]*common.Material)
	return nil
}

// SetLogLevel sets the log level
func SetLogLevel(level int) {
	log.SetLogLevel(level)
}

// Open is a smart opener for files quail might support without being within an archive
func Open(name string, r io.ReadSeeker) (interface{}, error) {
	var err error
	ext := strings.ToLower(filepath.Ext(name))
	name = filepath.Base(name)
	switch ext {
	case ".zon":
		zone := common.NewZone(name)
		err = zon.Decode(zone, r)
		if err != nil {
			return nil, fmt.Errorf("zone.Decode: %w", err)
		}
		return zone, nil
	case ".pts":
		point := common.NewParticlePoint(name)
		err = pts.Decode(point, r)
		if err != nil {
			return nil, fmt.Errorf("point.Decode: %w", err)
		}
		return point, nil
	case ".prt":
		particle := common.NewParticleRender(name)
		err = prt.Decode(particle, r)
		if err != nil {
			return nil, fmt.Errorf("particle.Decode: %w", err)
		}
		return particle, nil
	case ".lay":
		model := common.NewModel(name)
		err = lay.Decode(model, r)
		if err != nil {
			return nil, fmt.Errorf("model.Decode: %w", err)
		}
		return model, nil
	case ".lit":
		lits := []*common.RGBA{}
		err = lit.Decode(lits, r)
		if err != nil {
			return nil, fmt.Errorf("lit.Decode: %w", err)
		}
		return lits, nil
	case ".ani":
		animation := common.NewAnimation(name)
		err = ani.Decode(animation, r)
		if err != nil {
			return nil, fmt.Errorf("animation.Decode: %w", err)
		}
		return animation, nil
	case ".mod":
		model := common.NewModel(name)
		err = mod.Decode(model, r)
		if err != nil {
			return nil, fmt.Errorf("model.Decode: %w", err)
		}
		return model, nil
	case ".ter":
		model := common.NewModel(name)
		err = ter.Decode(model, r)
		if err != nil {
			return nil, fmt.Errorf("terrain.Decode: %w", err)
		}
		return model, nil
	case ".mds":
		model := common.NewModel(name)
		err = mds.Decode(model, r)
		if err != nil {
			return nil, fmt.Errorf("model.Decode: %w", err)
		}
		return model, nil
	case ".wld":
		q := New()
		wld, err := q.WLDDecode(r, nil)
		if err != nil {
			return nil, fmt.Errorf("wld.Decode: %w", err)
		}
		err = q.WldUnmarshal(wld)
		if err != nil {
			return nil, fmt.Errorf("wld.Unmarshal: %w", err)
		}
		return q, nil
	}

	return nil, fmt.Errorf("unknown extension %s", ext)
}
