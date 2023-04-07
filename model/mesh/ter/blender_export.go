package ter

import (
	"fmt"
	"os"

	"github.com/xackery/quail/dump"
)

// BlenderExport exports a TER file to a directory for use in blender
func (e *TER) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}

	pw, err := os.Create(fmt.Sprintf("%s/material_property.txt", path))
	if err != nil {
		return fmt.Errorf("create material_property.txt: %w", err)
	}
	defer pw.Close()
	pw.WriteString("material_name property_name value category\n")

	mw, err := os.Create(fmt.Sprintf("%s/material.txt", path))
	if err != nil {
		return fmt.Errorf("create material.txt: %w", err)
	}
	defer mw.Close()
	mw.WriteString("name flag shader_name\n")
	for _, m := range e.materials {
		mw.WriteString(m.Name + " ")
		mw.WriteString(dump.Str(m.Flag) + " ")
		mw.WriteString(m.ShaderName + "\n")
		for _, property := range m.Properties {
			pw.WriteString(m.Name + " ")
			pw.WriteString(property.Name + " ")
			pw.WriteString(dump.Str(property.Value) + " ")
			pw.WriteString(dump.Str(property.Category) + "\n")
		}
	}

	ppw, err := os.Create(fmt.Sprintf("%s/particle_point.txt", path))
	if err != nil {
		return fmt.Errorf("create particle_point.txt: %w", err)
	}
	defer ppw.Close()
	ppw.WriteString("name bone translation rotation scale\n")
	for _, pp := range e.particlePoints {
		ppw.WriteString(dump.Str(pp.Name) + " ")
		ppw.WriteString(dump.Str(pp.Bone) + " ")
		ppw.WriteString(dump.Str(pp.Translation) + " ")
		ppw.WriteString(dump.Str(pp.Rotation) + " ")
		ppw.WriteString(dump.Str(pp.Scale) + "\n")
	}

	prw, err := os.Create(fmt.Sprintf("%s/particle_render.txt", path))
	if err != nil {
		return fmt.Errorf("create particle_render.txt: %w", err)
	}
	defer prw.Close()
	prw.WriteString("duration id id2 particle_point unknownA unknownB unknownFFFFFFFF\n")
	for _, pr := range e.particleRenders {
		prw.WriteString(dump.Str(pr.Duration) + " ")
		prw.WriteString(dump.Str(pr.ID) + " ")
		prw.WriteString(dump.Str(pr.ID2) + " ")
		prw.WriteString(dump.Str(pr.ParticlePoint) + " ")
		prw.WriteString(dump.Str(pr.UnknownA) + " ")
		prw.WriteString(dump.Str(pr.UnknownB) + " ")
		prw.WriteString(dump.Str(pr.UnknownFFFFFFFF) + "\n")
	}

	tw, err := os.Create(fmt.Sprintf("%s/triangle.txt", path))
	if err != nil {
		return fmt.Errorf("create triangle.txt: %w", err)
	}
	defer tw.Close()
	tw.WriteString("index flag material_name\n")
	for _, t := range e.triangles {
		tw.WriteString(dump.Str(t.Index) + " ")
		tw.WriteString(dump.Str(t.Flag) + " ")
		tw.WriteString(dump.Str(t.MaterialName) + "\n")
	}

	vw, err := os.Create(fmt.Sprintf("%s/vertex.txt", path))
	if err != nil {
		return fmt.Errorf("create vertex.txt: %w", err)
	}
	defer vw.Close()
	vw.WriteString("joint position normal uv uv2 tint weight\n")
	for _, v := range e.vertices {
		vw.WriteString(dump.Str(v.Joint) + " ")
		vw.WriteString(dump.Str(v.Position) + " ")
		vw.WriteString(dump.Str(v.Normal) + " ")
		vw.WriteString(dump.Str(v.Uv) + " ")
		vw.WriteString(dump.Str(v.Uv2) + " ")
		vw.WriteString(dump.Str(v.Tint) + " ")
		vw.WriteString(dump.Str(v.Weight) + "\n")
	}

	return nil
}
