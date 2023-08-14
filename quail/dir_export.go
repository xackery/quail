package quail

import (
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/helper"
)

// DirExport exports the quail target to a directory
func (quail *Quail) DirExport(path string) error {

	path = strings.TrimSuffix(path, ".eqg")
	path = strings.TrimSuffix(path, ".s3d")
	path = strings.TrimSuffix(path, ".quail")
	path += ".quail"

	_, err := os.Stat(path)
	if err == nil {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("path %s is not a directory", path)
	}

	for _, mesh := range quail.Meshes {
		meshPath := fmt.Sprintf("%s/%s.mesh", path, mesh.Name)
		err = os.MkdirAll(meshPath, 0755)
		if err != nil {
			return fmt.Errorf("mkdir %s: %w", mesh.Name, err)
		}

		if len(mesh.Triangles) > 0 {
			tw, err := os.Create(fmt.Sprintf("%s/triangle.txt", meshPath))
			if err != nil {
				return fmt.Errorf("create mesh %s triangle.txt: %w", mesh.Name, err)
			}
			defer tw.Close()

			tw.WriteString("index|flag|material_name\n")
			tw.WriteString(fmt.Sprintf("ext|%s|-1\n", mesh.FileType))
			for _, triangle := range mesh.Triangles {
				tw.WriteString(fmt.Sprintf("%d,%d,%d|%d|%s\n", triangle.Index.X, triangle.Index.Y, triangle.Index.Z, triangle.Flag, triangle.MaterialName))
			}
		}

		if len(mesh.Vertices) > 0 {
			vw, err := os.Create(fmt.Sprintf("%s/vertex.txt", meshPath))
			if err != nil {
				return fmt.Errorf("create mesh %s vertex.txt: %w", mesh.Name, err)
			}
			defer vw.Close()

			vw.WriteString("position|normal|uv|uv2|tint\n")
			for _, vertex := range mesh.Vertices {

				vw.WriteString(fmt.Sprintf("%0.8f,%0.8f,%0.8f|%0.8f,%0.8f,%0.8f|%0.8f,%0.8f|%0.8f,%0.8f|%d,%d,%d,%d\n",
					vertex.Position.X, vertex.Position.Y, vertex.Position.Z,
					vertex.Normal.X, vertex.Normal.Y, vertex.Normal.Z,
					vertex.Uv.X, vertex.Uv.Y,
					vertex.Uv2.X, vertex.Uv2.Y,
					vertex.Tint.R, vertex.Tint.G, vertex.Tint.B, vertex.Tint.A))
			}
		}

		if len(mesh.Bones) > 0 {
			bw, err := os.Create(fmt.Sprintf("%s/bone.txt", meshPath))
			if err != nil {
				return fmt.Errorf("create mesh %s bone.txt: %w", mesh.Name, err)
			}
			defer bw.Close()

			bw.WriteString("name|child_index|children_count|next|pivot|rotation|scale\n")
			for _, bone := range mesh.Bones {
				bw.WriteString(fmt.Sprintf("%s|%d|%d|%d", bone.Name, bone.ChildIndex, bone.ChildrenCount, bone.Next))
				bw.WriteString(fmt.Sprintf("|%0.8f,%0.8f,%0.8f", bone.Pivot.Y, -bone.Pivot.X, bone.Pivot.Z)) //xyz is wonky
				bw.WriteString(fmt.Sprintf("|%0.8f,%0.8f,%0.8f,%0.8f", bone.Rotation.X, bone.Rotation.Y, bone.Rotation.Z, bone.Rotation.W))
				bw.WriteString(fmt.Sprintf("|%0.8f,%0.8f,%0.8f\n", bone.Scale.X, bone.Scale.Y, bone.Scale.Z))
			}
		}

		if len(mesh.ParticleRenders) > 0 {
			prw, err := os.Create(fmt.Sprintf("%s/particle_render.txt", meshPath))
			if err != nil {
				return fmt.Errorf("create mesh %s particle_render.txt: %w", mesh.Name, err)
			}
			defer prw.Close()

			prw.WriteString("id|id2|particle_point|unknowna1|unknowna2|unknowna3|unknowna4|unknowna5|duration|unknownb|unknownffffffff|unknownc\n")
			for _, render := range mesh.ParticleRenders {
				for _, entry := range render.Entries {
					prw.WriteString(fmt.Sprintf("%d|%d|%s|", entry.ID, entry.ID2, entry.ParticlePoint))
					prw.WriteString(fmt.Sprintf("%d|%d|%d|%d|%d|", entry.UnknownA1, entry.UnknownA2, entry.UnknownA3, entry.UnknownA4, entry.UnknownA5))
					prw.WriteString(fmt.Sprintf("%d|%d|%d|%d\n", entry.Duration, entry.UnknownB, entry.UnknownFFFFFFFF, entry.UnknownC))
				}
			}
		}

		if len(mesh.ParticlePoints) > 0 {
			ppw, err := os.Create(fmt.Sprintf("%s/particle_point.txt", meshPath))
			if err != nil {
				return fmt.Errorf("create mesh %s particle_point.txt: %w", mesh.Name, err)
			}
			defer ppw.Close()

			ppw.WriteString("name|bone|translation|rotation|scale\n")
			for _, point := range mesh.ParticlePoints {
				for _, entry := range point.Entries {
					ppw.WriteString(fmt.Sprintf("%s|%s|%0.8f,%0.8f,%0.8f|%0.8f,%0.8f,%0.8f|%0.8f,%0.8f,%0.8f\n", entry.Name, entry.Bone, entry.Translation.X, entry.Translation.Y, entry.Translation.Z, entry.Rotation.X, entry.Rotation.Y, entry.Rotation.Z, entry.Scale.X, entry.Scale.Y, entry.Scale.Z))
				}
			}
		}

		for _, material := range mesh.Materials {
			materialPath := fmt.Sprintf("%s/%s.material", path, material.Name)
			_, err = os.Stat(materialPath)
			if err == nil {
				continue
			}
			err = os.MkdirAll(materialPath, 0755)
			if err != nil {
				return err
			}

			mw, err := os.Create(fmt.Sprintf("%s/property.txt", materialPath))
			if err != nil {
				return fmt.Errorf("create mesh %s material %s property.txt: %w", mesh.Name, material.Name, err)
			}
			defer mw.Close()

			mw.WriteString("property_name|value|category\n")
			mw.WriteString(fmt.Sprintf("shaderName|%s|2\n", material.ShaderName))
			for _, property := range material.Properties {
				value := strings.ToLower(property.Value)
				if strings.ToLower(property.Name) == "e_fshininess0" {
					val := helper.AtoF32(property.Value)
					if val > 100 {
						val = 1.0
					} else {
						val /= 100
					}
					value = fmt.Sprintf("%0.8f", val)
				}
				mw.WriteString(fmt.Sprintf("%s|%s|%d\n", property.Name, value, property.Category))
				if len(property.Data) > 0 {
					err = os.WriteFile(fmt.Sprintf("%s/%s", materialPath, property.Value), property.Data, 0644)
					if err != nil {
						return err
					}
				}
			}

		}
	}

	for _, anim := range quail.Animations {

		animPath := fmt.Sprintf("%s/%s.ani", path, anim.Name)
		err = os.MkdirAll(animPath, 0755)
		if err != nil {
			return fmt.Errorf("mkdir %s: %w", anim.Name, err)
		}

		aw, err := os.Create(fmt.Sprintf("%s/animation.txt", animPath))
		if err != nil {
			return fmt.Errorf("create anim %s: %w", anim.Name, err)
		}
		defer aw.Close()

		aw.WriteString("key|value\n")

		val := 0
		if anim.IsStrict {
			val = 1
		}
		aw.WriteString(fmt.Sprintf("is_strict|%d\n", val))

		for _, bone := range anim.Bones {
			fw, err := os.Create(fmt.Sprintf("%s/%s.txt", animPath, bone.Name))
			if err != nil {
				return fmt.Errorf("create anim %s bone %s: %w", anim.Name, bone.Name, err)
			}
			defer fw.Close()

			fw.WriteString("milliseconds|rotation|scale|translation\n")
			for _, frame := range bone.Frames {
				fw.WriteString(fmt.Sprintf("%d|", frame.Milliseconds))
				fw.WriteString(fmt.Sprintf("%0.8f,%0.8f,%0.8f,%0.8f|", frame.Rotation.X, frame.Rotation.Y, frame.Rotation.Z, frame.Rotation.W))
				fw.WriteString(fmt.Sprintf("%0.8f,%0.8f,%0.8f|", frame.Scale.X, frame.Scale.Y, frame.Scale.Z))
				fw.WriteString(fmt.Sprintf("%0.8f,%0.8f,%0.8f\n", frame.Translation.X, frame.Translation.Y, frame.Translation.Z))
			}
		}
	}

	return nil
}
