package mds

import (
	"fmt"
	"os"

	"github.com/xackery/quail/dump"
)

// BlenderExport exports the MDS to a directory for use in Blender.
func (e *MDS) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	if len(e.animations) > 0 {
		aw, err := os.Create(fmt.Sprintf("%s/animation.txt", path))
		if err != nil {
			return fmt.Errorf("create animation.txt: %w", err)
		}
		defer aw.Close()
		aw.WriteString("name\n")
		for _, anim := range e.animations {
			aw.WriteString(anim.Name + "\n")
		}
	}

	if len(e.materials) > 0 {

		pw, err := os.Create(fmt.Sprintf("%s/material_property.txt", path))
		if err != nil {
			return fmt.Errorf("create material_property.txt: %w", err)
		}
		defer pw.Close()
		pw.WriteString("material_name|property_name|value|category\n")

		mw, err := os.Create(fmt.Sprintf("%s/material.txt", path))
		if err != nil {
			return fmt.Errorf("create material.txt: %w", err)
		}
		mw.WriteString("name|flag|shader_name\n")
		defer mw.Close()
		for _, m := range e.materials {
			mw.WriteString(m.Name + "|")
			mw.WriteString(dump.Str(m.Flag) + "|")
			mw.WriteString(m.ShaderName + "\n")
			for _, property := range m.Properties {
				pw.WriteString(m.Name + "|")
				pw.WriteString(property.Name + "|")
				pw.WriteString(dump.Str(property.Value) + "|")
				pw.WriteString(dump.Str(property.Category) + "\n")
			}
		}
	}

	if len(e.particlePoints) > 0 {

		ppw, err := os.Create(fmt.Sprintf("%s/particle_point.txt", path))
		if err != nil {
			return fmt.Errorf("create particle_point.txt: %w", err)
		}
		defer ppw.Close()
		ppw.WriteString("name|bone|translation|rotation|scale\n")
		for _, pp := range e.particlePoints {
			ppw.WriteString(dump.Str(pp.Name) + "|")
			ppw.WriteString(dump.Str(pp.Bone) + "|")
			ppw.WriteString(dump.Str(pp.Translation) + "|")
			ppw.WriteString(dump.Str(pp.Rotation) + "|")
			ppw.WriteString(dump.Str(pp.Scale) + "\n")
		}
	}

	if len(e.particleRenders) > 0 {
		prw, err := os.Create(fmt.Sprintf("%s/particle_render.txt", path))
		if err != nil {
			return fmt.Errorf("create particle_render.txt: %w", err)
		}
		defer prw.Close()
		prw.WriteString("duration|id|id2|particle_point|unknownA|unknownB|unknownFFFFFFFF\n")
		for _, pr := range e.particleRenders {
			prw.WriteString(dump.Str(pr.Duration) + "|")
			prw.WriteString(dump.Str(pr.ID) + "|")
			prw.WriteString(dump.Str(pr.ID2) + "|")
			prw.WriteString(dump.Str(pr.ParticlePoint) + "|")
			prw.WriteString(dump.Str(pr.UnknownA) + "|")
			prw.WriteString(dump.Str(pr.UnknownB) + "|")
			prw.WriteString(dump.Str(pr.UnknownFFFFFFFF) + "\n")
		}
	}

	if len(e.bones) > 0 {
		sw, err := os.Create(fmt.Sprintf("%s/skin.txt", path))
		if err != nil {
			return fmt.Errorf("create skin.txt: %w", err)
		}
		defer sw.Close()

		bw, err := os.Create(fmt.Sprintf("%s/bone.txt", path))
		if err != nil {
			return fmt.Errorf("create bone.txt: %w", err)
		}
		defer bw.Close()
		bw.WriteString("name|child_index|children_count|next|pivot|rotation|scale\n")
		for _, b := range e.bones {
			bw.WriteString(dump.Str(b.Name) + "|")
			bw.WriteString(dump.Str(b.ChildIndex) + "|")
			bw.WriteString(dump.Str(b.ChildrenCount) + "|")
			bw.WriteString(dump.Str(b.Next) + "|")
			bw.WriteString(dump.Str(b.Pivot) + "|")
			bw.WriteString(dump.Str(b.Rotation) + "|")
			bw.WriteString(dump.Str(b.Scale) + "\n")
		}
	}

	/*for _, file := range e.files {
		fw, err := os.Create(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil {
			return fmt.Errorf("create %s: %w", file.Name(), err)
		}
		defer fw.Close()
		fw.Write(file.Data())
	}
	*/

	if len(e.triangles) > 0 {
		tw, err := os.Create(fmt.Sprintf("%s/triangle.txt", path))
		if err != nil {
			return fmt.Errorf("create triangle.txt: %w", err)
		}
		defer tw.Close()
		tw.WriteString("index flag material_name\n")
		for _, t := range e.triangles {
			tw.WriteString(dump.Str(t.Index) + "|")
			tw.WriteString(dump.Str(t.Flag) + "|")
			tw.WriteString(dump.Str(t.MaterialName) + "\n")
		}

		vw, err := os.Create(fmt.Sprintf("%s/vertex.txt", path))
		if err != nil {
			return fmt.Errorf("create vertex.txt: %w", err)
		}
		defer vw.Close()
		vw.WriteString("joint|position|normal|uv|uv2|tint|weight\n")
		for _, v := range e.vertices {
			vw.WriteString(dump.Str(v.Joint) + "|")
			vw.WriteString(dump.Str(v.Position) + "|")
			vw.WriteString(dump.Str(v.Normal) + "|")
			vw.WriteString(dump.Str(v.Uv) + "|")
			vw.WriteString(dump.Str(v.Uv2) + "|")
			vw.WriteString(dump.Str(v.Tint) + "|")
			vw.WriteString(dump.Str(v.Weight) + "\n")
		}
	}
	return nil
}
