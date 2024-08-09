package quail

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wld"
)

type Quail struct {
	Header                 *common.Header
	Models                 []*common.Model
	Animations             []*common.Animation
	Zone                   *common.Zone
	materialCache          map[string]*common.Material
	IsExtensionVersionDump bool
	Textures               map[string][]byte // Textures are raw texture files
	wld                    *wld.Wld
	wldObject              *wld.Wld
	wldLights              *wld.Wld
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
	//name = filepath.Base(name)
	switch ext {
	case ".zon":
		zon := &raw.Zon{}
		err = zon.Read(r)
		if err != nil {
			return nil, fmt.Errorf("zone.Decode: %w", err)
		}
		return zon, nil
	case ".pts":
		pts := &raw.Pts{}
		err = pts.Read(r)
		if err != nil {
			return nil, fmt.Errorf("point.Decode: %w", err)
		}
		return pts, nil
	case ".prt":
		prt := &raw.Prt{}
		err = prt.Read(r)
		if err != nil {
			return nil, fmt.Errorf("particle.Decode: %w", err)
		}
		return prt, nil
	case ".lay":
		lay := &raw.Lay{}
		err = lay.Read(r)
		if err != nil {
			return nil, fmt.Errorf("lay.Decode: %w", err)
		}
		return lay, nil
	case ".lit":
		lit := &raw.Lit{}
		err = lit.Read(r)
		if err != nil {
			return nil, fmt.Errorf("lit.Decode: %w", err)
		}
		return lit, nil
	case ".ani":
		ani := &raw.Ani{}
		err = ani.Read(r)
		if err != nil {
			return nil, fmt.Errorf("ani.Read: %w", err)
		}
		return ani, nil
	case ".mod":
		mod := &raw.Mod{}
		err = mod.Read(r)
		if err != nil {
			return nil, fmt.Errorf("mod.Read: %w", err)
		}
		return mod, nil
	case ".ter":
		ter := &raw.Ter{}
		err = ter.Read(r)
		if err != nil {
			return nil, fmt.Errorf("terrain.Decode: %w", err)
		}
		return ter, nil
	case ".mds":
		mds := &raw.Mds{}
		err = mds.Read(r)
		if err != nil {
			return nil, fmt.Errorf("mds.Decode: %w", err)
		}
		return mds, nil
	case ".wld":
		header := make([]byte, 4)
		_, err = r.Read(header)
		if err != nil {
			return nil, fmt.Errorf("read header: %w", err)
		}
		r.Seek(0, io.SeekStart)
		if string(header) != "\x02\x3D\x50\x54" {
			wldAscii := &raw.WldAscii{}
			err = wldAscii.Read(r)
			if err != nil {
				return nil, fmt.Errorf("wldAscii.Decode: %w", err)
			}
			return wldAscii, nil
		}

		wld := &raw.Wld{}
		err = wld.Read(r)
		if err != nil {
			return nil, fmt.Errorf("wld.Decode: %w", err)
		}
		return wld, nil
	case ".sph":
		return nil, nil
	case ".sps": // map file, safely ignored
		return nil, nil
	}

	return nil, fmt.Errorf("unknown extension %s", ext)
}
