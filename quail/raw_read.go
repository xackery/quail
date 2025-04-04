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
	case *raw.Wld:
		return q.wldRead(val, in.FileName())
	case *raw.Dds, *raw.Bmp, *raw.Png, *raw.Tga, *raw.Jpg: // textures
		return q.assetRead(val)
	case *raw.Lit: // baked lighting in eqg
		return q.assetRead(val)
	case *raw.Txt:
		return q.assetRead(val)
	case *raw.Mod, *raw.Pts, *raw.Prt, *raw.Mds, *raw.Ter, *raw.Lod, *raw.Lay, *raw.Ani, *raw.Tog, *raw.Zon, *raw.Dat, *raw.Eco, *raw.Rfd, *raw.Def:
		//fmt.Println("ignoring", in.Identity())
		return nil // ignored, loaded by wce parsre
	case *raw.Unk:
		return q.assetRead(val)
	default:
		return fmt.Errorf("unknown type %T (%s)", val, val)
	}
}

func RawRead(in raw.ReadWriter, q *Quail) error {
	if q == nil {
		return fmt.Errorf("quail is nil")
	}
	return q.RawRead(in)
}

func (q *Quail) assetRead(in raw.ReadWriter) error {
	buf := &bytes.Buffer{}
	err := in.Write(buf)
	if err != nil {
		return fmt.Errorf("write asset %s: %w", in.Identity(), err)
	}
	q.assetAdd(in.FileName(), buf.Bytes())
	return nil
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
