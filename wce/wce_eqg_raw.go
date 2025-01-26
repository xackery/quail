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
	case ".lay":
		rawSrc := &raw.Lay{
			MetaFileName: strings.TrimSuffix(entry.Name(), ".lay"),
		}
		err = rawSrc.Read(bytes.NewReader(entry.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		def := &LayDef{}
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

		dst.Materials, err = writeEqgMaterials(mds.Materials)
		if err != nil {
			return fmt.Errorf("write materials: %w", err)
		}

		for _, bone := range mds.Bones {
			dst.Bones = append(dst.Bones, &raw.Bone{
				Name:          bone.Name,
				Next:          bone.Next,
				ChildrenCount: bone.ChildrenCount,
				ChildIndex:    bone.ChildIndex,
				Pivot:         bone.Pivot,
				Quaternion:    bone.Quaternion,
				Scale:         bone.Scale,
			})
		}

		for _, model := range mds.Models {
			model := &raw.MdsModel{
				MainPiece: model.MainPiece,
				Name:      model.Name,
			}
			for _, vert := range model.Vertices {
				model.Vertices = append(model.Vertices, &raw.Vertex{
					Position: vert.Position,
					Normal:   vert.Normal,
					Tint:     vert.Tint,
					Uv:       vert.Uv,
					Uv2:      vert.Uv2,
				})
			}
			for _, face := range model.Faces {
				flags := uint32(0)
				if face.Flags&1 == 1 {
					flags |= 1
				}

				model.Faces = append(model.Faces, &raw.Face{
					Index:        face.Index,
					MaterialName: face.MaterialName,
					Flags:        flags,
				})

			}

			for _, boneAssignment := range model.BoneAssignments {
				weights := [4]*raw.MdsBoneWeight{}
				for i := 0; i < len(boneAssignment); i++ {
					wt := boneAssignment[i]
					weight := &raw.MdsBoneWeight{
						BoneIndex: wt.BoneIndex,
						Value:     wt.Value,
					}
					weights[i] = weight
				}
				model.BoneAssignments = append(model.BoneAssignments, weights)
			}

			dst.Models = append(dst.Models, model)
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

		dst.Materials, err = writeEqgMaterials(mod.Materials)
		if err != nil {
			return fmt.Errorf("write materials: %w", err)
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
		dst := &raw.Ter{}

		materialNames := make(map[string]bool)
		for _, face := range ter.Faces {
			materialNames[face.MaterialName] = true
		}

		materials := []*EQMaterialDef{}
		for materialName := range materialNames {
			isFound := false
			for _, material := range wce.EQMaterialDefs {
				if material.Tag != materialName {
					continue
				}
				materials = append(materials, material)
				isFound = true
				break
			}
			if !isFound {
				return fmt.Errorf("terrain %s refers to material %s, but not declared", ter.Tag, materialName)
			}
		}

		dst.Materials, err = writeEqgMaterials(materials)
		if err != nil {
			return fmt.Errorf("write materials: %w", err)
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

func writeEqgMaterials(srcMaterials []*EQMaterialDef) ([]*raw.Material, error) {
	dstMaterials := []*raw.Material{}
	for _, srcMat := range srcMaterials {
		mat := &raw.Material{
			Name:       srcMat.Tag,
			EffectName: srcMat.ShaderTag,
		}
		if srcMat.HexOneFlag == 1 {
			mat.Flags &= 1
		}

		for _, prop := range srcMat.Properties {
			mat.Properties = append(mat.Properties, &raw.MaterialParam{
				Name:  prop.Name,
				Value: prop.Value,
				Type:  prop.Type,
			})
		}

		dstMaterials = append(dstMaterials, mat)
	}
	return dstMaterials, nil
}
