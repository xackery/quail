package wce

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func (wce *Wce) ReadEqgRaw(archive *pfs.Pfs) error {
	if archive == nil {
		return fmt.Errorf("archive is nil")
	}

	if !wce.WorldDef.EqgVersion.Valid {
		wce.WorldDef.EqgVersion.Valid = true
		wce.WorldDef.EqgVersion.Int8 = 1
	}
	baseName := archive.Name()
	ext := filepath.Ext(baseName)
	if ext != ".eqg" {
		return fmt.Errorf("invalid eqg file %s (no .eqg suffix)", baseName)
	}
	baseName = strings.TrimSuffix(baseName, ext)

	files := archive.Files()
	for _, file := range files {
		if strings.EqualFold(file.Name(), baseName+".dat") {
			wce.WorldDef.Zone = 1
			break
		}
		if strings.Contains(file.Name(), ".ter") {
			wce.WorldDef.Zone = 1
			break
		}
	}

	for _, file := range files {
		err := wce.readEqgEntry(file)
		if err != nil {
			return fmt.Errorf("read eqg entry %s: %w", file.Name(), err)
		}
	}

	return nil
}

func (wce *Wce) readEqgEntry(entry *pfs.FileEntry) error {
	var err error

	ext := strings.ToLower(filepath.Ext(entry.Name()))
	switch ext {
	case ".mds":
		rawSrc := &raw.Mds{
			MetaFileName: strings.TrimSuffix(entry.Name(), ".mds"),
		}
		err = rawSrc.Read(bytes.NewReader(entry.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		def := &MdsDef{}
		err := def.FromRaw(wce, rawSrc)
		if err != nil {
			return fmt.Errorf("mds: %w", err)
		}
		wce.MdsDefs = append(wce.MdsDefs, def)

	case ".mod":
		rawSrc := &raw.Mod{
			MetaFileName: strings.TrimSuffix(entry.Name(), ".mod"),
		}
		err = rawSrc.Read(bytes.NewReader(entry.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		def := &ModDef{}
		err := def.FromRaw(wce, rawSrc)
		if err != nil {
			return fmt.Errorf("mod: %w", err)
		}
		wce.ModDefs = append(wce.ModDefs, def)
	case ".ter":
		rawSrc := &raw.Ter{
			MetaFileName: strings.TrimSuffix(entry.Name(), ".ter"),
		}
		err = rawSrc.Read(bytes.NewReader(entry.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		def := &TerDef{}
		err := def.FromRaw(wce, rawSrc)
		if err != nil {
			return fmt.Errorf("ter: %w", err)
		}
		wce.TerDefs = append(wce.TerDefs, def)

	case ".ani":
		rawSrc := &raw.Ani{
			MetaFileName: strings.TrimSuffix(entry.Name(), ".ani"),
		}
		err = rawSrc.Read(bytes.NewReader(entry.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		def := &AniDef{}
		err := def.FromRaw(wce, rawSrc)
		if err != nil {
			return fmt.Errorf("ani: %w", err)
		}
		wce.AniDefs = append(wce.AniDefs, def)
	default:
		return nil
	}

	return nil
}

func (wce *Wce) WriteEqgRaw(archive *pfs.Pfs) error {
	if archive == nil {
		return fmt.Errorf("archive is nil")
	}

	for _, mds := range wce.MdsDefs {
		buf := &bytes.Buffer{}
		dst := &raw.Mds{}
		err := dst.Write(buf)
		if err != nil {
			return fmt.Errorf("mds write: %w", err)
		}
		err = archive.Add(mds.Tag+".mds", buf.Bytes())
		if err != nil {
			return fmt.Errorf("add mds: %w", err)
		}
	}

	for _, mod := range wce.ModDefs {
		buf := &bytes.Buffer{}
		dst := &raw.Mod{}
		err := dst.Write(buf)
		if err != nil {
			return fmt.Errorf("mod write: %w", err)
		}
		err = archive.Add(mod.Tag+".mod", buf.Bytes())
		if err != nil {
			return fmt.Errorf("add mod: %w", err)
		}
	}

	for _, ter := range wce.TerDefs {
		buf := &bytes.Buffer{}
		dst := &raw.Ter{}
		err := dst.Write(buf)
		if err != nil {
			return fmt.Errorf("ter write: %w", err)
		}
		err = archive.Add(ter.Tag+".ter", buf.Bytes())
		if err != nil {
			return fmt.Errorf("add ter: %w", err)
		}
	}

	for _, ani := range wce.AniDefs {
		buf := &bytes.Buffer{}
		dst := &raw.Ani{}
		err := dst.Write(buf)
		if err != nil {
			return fmt.Errorf("ani write: %w", err)
		}
		err = archive.Add(ani.Tag+".ani", buf.Bytes())
		if err != nil {
			return fmt.Errorf("add ani: %w", err)
		}
	}

	return nil
}
