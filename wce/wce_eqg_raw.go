package wce

import (
	"fmt"
	"io"

	"github.com/xackery/quail/raw"
)

func (wce *Wce) ReadEqgRaw(src raw.Reader) error {

	if !wce.WorldDef.EqgVersion.Valid {
		wce.WorldDef.EqgVersion.Valid = true
		wce.WorldDef.EqgVersion.Int8 = 1
	}
	switch frag := src.(type) {
	case *raw.Ani:
		def := &AniDef{}
		err := def.FromRaw(wce, frag)
		if err != nil {
			return fmt.Errorf("ani: %w", err)
		}
		wce.AniDefs = append(wce.AniDefs, def)
	case *raw.Mds:
		def := &MdsDef{}
		err := def.FromRaw(wce, frag)
		if err != nil {
			return fmt.Errorf("mds: %w", err)
		}
		wce.MdsDefs = append(wce.MdsDefs, def)
	case *raw.Mod:
		def := &ModDef{}
		err := def.FromRaw(wce, frag)
		if err != nil {
			return fmt.Errorf("mod: %w", err)
		}
		wce.ModDefs = append(wce.ModDefs, def)
	case *raw.Ter:
		def := &TerDef{}
		err := def.FromRaw(wce, frag)
		if err != nil {
			return fmt.Errorf("ter: %w", err)
		}
		wce.TerDefs = append(wce.TerDefs, def)
	}

	return nil
}

func (wce *Wce) WriteEqgRaw(w io.Writer) error {
	return nil
}
