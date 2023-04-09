package wld

import (
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/dump"
	"github.com/xackery/quail/model/geo"
)

// BlenderExport exports WLD to a blender dir
func (e *WLD) BlenderExport(dir string) error {
	path := fmt.Sprintf("%s/_%s", dir, e.Name())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
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
				val := property.Name
				if strings.EqualFold(property.Value, "e_texturediffuse0") {
					val = "e_TextureDiffuse0"
				}
				pw.WriteString(val + "|")
				pw.WriteString(dump.Str(property.Value) + "|")
				pw.WriteString(dump.Str(property.Category) + "\n")
			}
		}
	}

	if len(e.meshes) > 0 {

		for i, mesh := range e.meshes {
			fmt.Println("exporting mesh", mesh.Name)
			if mesh.Name == "" {
				mesh.Name = fmt.Sprintf("mesh_%d", i)
			}
			err = e.blenderExportMesh(dir, mesh)
			if err != nil {
				return fmt.Errorf("blenderExportMeshes: %w", err)
			}
		}
	}

	return nil
}

func (e *WLD) blenderExportMesh(dir string, mesh *geo.Mesh) error {
	if len(e.meshes) == 0 {
		return nil
	}

	path := fmt.Sprintf("%s/_%s.mds", dir, mesh.Name)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
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
				val := property.Name
				if strings.EqualFold(property.Value, "e_texturediffuse0") {
					val = "e_TextureDiffuse0"
				}
				pw.WriteString(val + "|")
				pw.WriteString(dump.Str(property.Value) + "|")
				pw.WriteString(dump.Str(property.Category) + "\n")
			}
		}
	}

	tw, err := os.Create(fmt.Sprintf("%s/triangle.txt", path))
	if err != nil {
		return fmt.Errorf("create triangle.txt: %w", err)
	}
	defer tw.Close()
	tw.WriteString("index flag material_name\n")
	for _, t := range mesh.Triangles {
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
	for _, v := range mesh.Vertices {
		vw.WriteString(dump.Str(v.Joint) + "|")
		vw.WriteString(dump.Str(v.Position) + "|")
		vw.WriteString(dump.Str(v.Normal) + "|")
		vw.WriteString(dump.Str(v.Uv) + "|")
		vw.WriteString(dump.Str(v.Uv2) + "|")
		vw.WriteString(dump.Str(v.Tint) + "|")
		vw.WriteString(dump.Str(v.Weight) + "\n")
	}
	return nil
}
