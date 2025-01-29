package quail

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/wce"
)

// Read takes a raw type and converts it to a quail type
func (q *Quail) RawRead(in raw.ReadWriter) error {
	if q == nil {
		return fmt.Errorf("quail is nil")
	}
	switch val := in.(type) {
	case *raw.Lay:
		return q.layRead(val)
	case *raw.Ani:
		return q.aniRead(val)
	case *raw.Wld:
		return q.wldRead(val, in.FileName())
	case *raw.Dds:
		return q.ddsRead(val)
	case *raw.Bmp:
		return q.bmpRead(val)
	case *raw.Png:
		return q.pngRead(val)
	case *raw.Mod:
		return q.modRead(val)
	case *raw.Pts:
		return nil
	case *raw.Prt:
		return nil
	case *raw.Mds:
		return q.mdsRead(val)
	case *raw.Unk:
		return nil
	case *raw.Txt:
		return nil
	default:
		return fmt.Errorf("unknown type %T", val)
	}
}

func RawRead(in raw.ReadWriter, q *Quail) error {
	if q == nil {
		return fmt.Errorf("quail is nil")
	}
	return q.RawRead(in)
}

func (q *Quail) aniRead(in *raw.Ani) error {
	return fmt.Errorf("ani not implemented")

}

func (q *Quail) layRead(in *raw.Lay) error {
	return fmt.Errorf("lay not implemented")
}

func (q *Quail) ddsRead(in *raw.Dds) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write dds: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) bmpRead(in *raw.Bmp) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write bmp: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) pngRead(in *raw.Png) error {
	if q.Textures == nil {
		q.Textures = make(map[string][]byte)
	}
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write png: %w", err)
	}
	q.Textures[in.FileName()] = buf.Bytes()
	return nil
}

func (q *Quail) modRead(in *raw.Mod) error {
	return fmt.Errorf("mod not implemented")
}

func (q *Quail) mdsRead(in *raw.Mds) error {
	return fmt.Errorf("mds not implemented")
}

func (q *Quail) wldRead(srcWld *raw.Wld, filename string) error {

	wld := wce.New(filename)
	err := wld.ReadWldRaw(srcWld)
	if err != nil {
		return fmt.Errorf("read wld: %w", err)
	}

	if strings.ToLower(filename) == "objects.wld" {
		q.WldObject = wld
	} else if strings.ToLower(filename) == "lights.wld" {
		q.WldLights = wld
	} else {
		q.Wld = wld
	}

	return nil
}
