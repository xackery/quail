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
		def := &EqgMdsDef{}
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
		def := &EqgModDef{}
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
		def := &EqgTerDef{}
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
		def := &EqgAniDef{}
		err := def.FromRaw(wce, rawSrc)
		if err != nil {
			return fmt.Errorf("ani: %w", err)
		}
		wce.AniDefs = append(wce.AniDefs, def)
	case ".lay":
		rawSrc := &raw.Lay{
			MetaFileName: strings.TrimSuffix(entry.Name(), ".lay"),
		}
		err = rawSrc.Read(bytes.NewReader(entry.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		def := &EqgLayDef{}
		err := def.FromRaw(wce, rawSrc)
		if err != nil {
			return fmt.Errorf("lay: %w", err)
		}
		wce.LayDefs = append(wce.LayDefs, def)
	default:
		return nil
	}

	return nil
}

func (wce *Wce) WriteEqgRaw(archive *pfs.Pfs) error {
	if archive == nil {
		return fmt.Errorf("archive is nil")
	}
	var err error

	for _, mds := range wce.MdsDefs {
		buf := &bytes.Buffer{}
		dst := &raw.Mds{
			MetaFileName: mds.Tag,
			Version:      mds.Version,
		}

		err = mds.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("mds to raw: %w", err)
		}

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
		dst := &raw.Mod{
			MetaFileName: mod.Tag,
			Version:      mod.Version,
		}

		err = mod.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("mod to raw: %w", err)
		}

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
		dst := &raw.Ter{
			MetaFileName: ter.Tag,
			Version:      ter.Version,
		}

		err = ter.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("ter to raw: %w", err)
		}

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
		err = ani.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("ani to raw: %w", err)
		}

		err = dst.Write(buf)
		if err != nil {
			return fmt.Errorf("ani write: %w", err)
		}
		err = archive.Add(ani.Tag+".ani", buf.Bytes())
		if err != nil {
			return fmt.Errorf("add ani: %w", err)
		}
	}

	for _, lay := range wce.LayDefs {
		buf := &bytes.Buffer{}
		dst := &raw.Lay{
			MetaFileName: lay.Tag,
			Version:      lay.Version,
		}
		for _, layEntry := range lay.Layers {
			dst.Layers = append(dst.Layers, &raw.LayEntry{
				Material: layEntry.Material,
				Diffuse:  layEntry.Diffuse,
				Normal:   layEntry.Normal,
			})
		}

		err = lay.ToRaw(wce, dst)
		if err != nil {
			return fmt.Errorf("lay to raw: %w", err)
		}

		err := dst.Write(buf)
		if err != nil {
			return fmt.Errorf("lay write: %w", err)
		}
		err = archive.Add(lay.Tag+".lay", buf.Bytes())
		if err != nil {
			return fmt.Errorf("add lay: %w", err)
		}
	}

	return nil
}

func writeEqgMaterials(srcMaterials []*EQMaterialDef) ([]*raw.ModMaterial, error) {
	dstMaterials := []*raw.ModMaterial{}
	for _, srcMat := range srcMaterials {
		mat := &raw.ModMaterial{
			Name:       srcMat.Tag,
			EffectName: srcMat.ShaderTag,
		}
		if srcMat.HexOneFlag == 1 {
			mat.Flags &= 1
		}

		for _, prop := range srcMat.Properties {
			mat.Properties = append(mat.Properties, &raw.ModMaterialParam{
				Name:  prop.Name,
				Value: prop.Value,
				Type:  prop.Type,
			})
		}

		dstMaterials = append(dstMaterials, mat)
	}
	return dstMaterials, nil
}
