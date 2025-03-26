package quail

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/qfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

type Quail struct {
	IsExtensionVersionDump bool
	Wld                    *wce.Wce
	WldObject              *wce.Wce
	WldLights              *wce.Wce
	Assets                 map[string][]byte
	StatFS                 fs.StatFS
	ReadDirFS              fs.ReadDirFS
	FileSystem             qfs.QFS
}

// New returns a new Quail instance
func New() *Quail {
	return &Quail{}
}

// Close flushes any memory and closes any open files
func (e *Quail) Close() error {
	return nil
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
	case ".lod":
		lod := &raw.Lod{}
		err = lod.Read(r)
		if err != nil {
			return nil, fmt.Errorf("lod.Decode: %w", err)
		}
		return lod, nil
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
		_, err = r.Seek(0, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("seek start: %w", err)
		}
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
	case ".db":
		return nil, nil
	case ".uvw":
		return nil, nil
	case ".max":
		return nil, nil
	case ".sph":
		return nil, nil
	case ".sps": // map file, safely ignored
		return nil, nil
	}

	return nil, fmt.Errorf("unknown extension %s", ext)
}
