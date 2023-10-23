package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// DirImport imports the quail target from a directory
func (quail *Quail) DirImport(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("path %s is not a directory", path)
	}
	if filepath.Ext(path) != ".quail" {
		return fmt.Errorf("path %s is not a .quail target", path)
	}

	quailFiles, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", path, err)
	}

	for _, qf := range quailFiles {
		if !qf.IsDir() {
			continue
		}

		switch filepath.Ext(qf.Name()) {
		case ".material":
			err = quail.dirParseMaterial(path, qf.Name())
			if err != nil {
				return fmt.Errorf("parse material %s: %w", qf.Name(), err)
			}
		}
	}
	for _, qf := range quailFiles {
		if !qf.IsDir() {
			continue
		}
		switch filepath.Ext(qf.Name()) {
		case ".model":
			err = quail.dirParseModel(path, qf.Name())
			if err != nil {
				return fmt.Errorf("parse model %s: %w", qf.Name(), err)
			}
		case ".ani":
			err = quail.dirParseAni(path, qf.Name())
			if err != nil {
				return fmt.Errorf("parse ani %s: %w", qf.Name(), err)
			}
		}
	}
	return nil
}

func (quail *Quail) dirParseModel(path string, name string) error {
	model := common.NewModel("")
	model.FileType = "mod" // default to mod
	model.Header.Name = strings.TrimSuffix(name, ".model")
	modelPath := fmt.Sprintf("%s/%s", path, name)
	modelFiles, err := os.ReadDir(modelPath)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", modelPath, err)
	}
	for _, mf := range modelFiles {
		if mf.Name() == "triangle.txt" {
			lines, err := helper.ReadFile(fmt.Sprintf("%s/%s", modelPath, mf.Name()))
			if err != nil {
				return fmt.Errorf("read model %s triangle.txt: %w", model.Header.Name, err)
			}
			for i, line := range lines {
				if i == 0 {
					continue
				}
				if len(line) == 0 {
					continue
				}
				records := strings.Split(line, "|")
				if len(records) != 3 {
					return fmt.Errorf("triangle.txt line %d: expected 3 records, got %d", i, len(records))
				}
				if records[0] == "ext" {
					model.FileType = records[1]
					continue
				}
				triangle := common.Triangle{}
				vec3 := strings.Split(records[0], ",")
				triangle.Index.X = helper.ParseUint32(vec3[0], 0)
				triangle.Index.Y = helper.ParseUint32(vec3[1], 0)
				triangle.Index.Z = helper.ParseUint32(vec3[2], 0)
				triangle.Flag = helper.ParseUint32(records[1], 0)
				triangle.MaterialName = helper.Clean(strings.TrimSuffix(strings.TrimSuffix(records[2], "\n"), "\r"))
				isLoaded := false
				for _, material := range model.Materials {
					if material.Name != triangle.MaterialName {
						continue
					}
					isLoaded = true
					break
				}
				if !isLoaded {
					mat, ok := quail.materialCache[triangle.MaterialName]
					if !ok {
						return fmt.Errorf("material %s not found", triangle.MaterialName)
					}
					model.Materials = append(model.Materials, mat)
				}
				model.Triangles = append(model.Triangles, triangle)
			}
		}
		if mf.Name() == "vertex.txt" {
			lines, err := helper.ReadFile(fmt.Sprintf("%s/%s", modelPath, mf.Name()))
			if err != nil {
				return fmt.Errorf("read model %s vertex.txt: %w", model.Header.Name, err)
			}
			for i, line := range lines {
				if i == 0 {
					continue
				}
				if len(line) == 0 {
					continue
				}
				records := strings.Split(line, "|")
				if len(records) != 5 {
					return fmt.Errorf("vertex.txt line %d: expected 5 records, got %d", i, len(records))
				}
				vertex := common.Vertex{}
				vec3 := strings.Split(records[0], ",")
				vertex.Position.X = helper.ParseFloat32(vec3[0], 0)
				vertex.Position.Y = helper.ParseFloat32(vec3[1], 0)
				vertex.Position.Z = helper.ParseFloat32(vec3[2], 0)
				vec3 = strings.Split(records[1], ",")
				vertex.Normal.X = helper.ParseFloat32(vec3[0], 0)
				vertex.Normal.Y = helper.ParseFloat32(vec3[1], 0)
				vertex.Normal.Z = helper.ParseFloat32(vec3[2], 0)
				vec2 := strings.Split(records[2], ",")
				vertex.Uv.X = helper.ParseFloat32(vec2[0], 0)
				vertex.Uv.Y = helper.ParseFloat32(vec2[1], 0)
				vec2 = strings.Split(records[3], ",")
				vertex.Uv2.X = helper.ParseFloat32(vec2[0], 0)
				vertex.Uv2.Y = helper.ParseFloat32(vec2[1], 0)
				rgb4 := strings.Split(records[4], ",")
				vertex.Tint.R = helper.ParseUint8(rgb4[0], 0)
				vertex.Tint.G = helper.ParseUint8(rgb4[1], 0)
				vertex.Tint.B = helper.ParseUint8(rgb4[2], 0)
				vertex.Tint.A = helper.ParseUint8(rgb4[3], 0)

				model.Vertices = append(model.Vertices, vertex)
			}
		}

		if mf.Name() == "bone.txt" {
			lines, err := helper.ReadFile(fmt.Sprintf("%s/%s", modelPath, mf.Name()))
			if err != nil {
				return fmt.Errorf("read model %s bone.txt: %w", model.Header.Name, err)
			}
			for i, line := range lines {
				if i == 0 {
					continue
				}
				if len(line) == 0 {
					continue
				}
				records := strings.Split(line, "|")
				if len(records) != 7 {
					return fmt.Errorf("bone.txt line %d: expected 7 records, got %d", i, len(records))
				}
				bone := common.Bone{}
				bone.Name = records[0]
				bone.ChildIndex = helper.ParseInt32(records[1], 0)
				bone.ChildrenCount = helper.ParseUint32(records[2], 0)
				bone.Next = helper.ParseInt32(records[3], 0)
				vec3 := strings.Split(records[4], ",")
				bone.Pivot.X = helper.ParseFloat32(vec3[0], 0)
				bone.Pivot.Y = helper.ParseFloat32(vec3[1], 0)
				bone.Pivot.Z = helper.ParseFloat32(vec3[2], 0)
				vec4 := strings.Split(records[5], ",")
				bone.Rotation.X = helper.ParseFloat32(vec4[0], 0)
				bone.Rotation.Y = helper.ParseFloat32(vec4[1], 0)
				bone.Rotation.Z = helper.ParseFloat32(vec4[2], 0)
				bone.Rotation.W = helper.ParseFloat32(vec4[3], 0)
				vec3 = strings.Split(records[6], ",")
				bone.Scale.X = helper.ParseFloat32(vec3[0], 0)
				bone.Scale.Y = helper.ParseFloat32(vec3[1], 0)
				bone.Scale.Z = helper.ParseFloat32(vec3[2], 0)

				model.Bones = append(model.Bones, bone)
			}
		}

		if mf.Name() == "particle_render.txt" {
			particleRender := common.NewParticleRender("")

			lines, err := helper.ReadFile(fmt.Sprintf("%s/%s", modelPath, mf.Name()))
			if err != nil {
				return fmt.Errorf("read model %s particle_render.txt: %w", model.Header.Name, err)
			}
			for i, line := range lines {
				if i == 0 {
					continue
				}
				if len(line) == 0 {
					continue
				}
				records := strings.Split(line, "|")
				if len(records) != 12 {
					return fmt.Errorf("particle_render.txt line %d: expected 11 records, got %d", i, len(records))
				}

				entry := &common.ParticleRenderEntry{}
				entry.ID = helper.ParseUint32(records[0], 0)
				entry.ID2 = helper.ParseUint32(records[1], 0)
				entry.ParticlePoint = records[2]
				entry.UnknownA1 = helper.ParseUint32(records[3], 0)
				entry.UnknownA2 = helper.ParseUint32(records[4], 0)
				entry.UnknownA3 = helper.ParseUint32(records[5], 0)
				entry.UnknownA4 = helper.ParseUint32(records[6], 0)
				entry.UnknownA5 = helper.ParseUint32(records[7], 0)

				entry.Duration = helper.ParseUint32(records[8], 0)
				entry.UnknownB = helper.ParseUint32(records[9], 0)
				entry.UnknownFFFFFFFF = helper.ParseInt32(records[10], 0)
				entry.UnknownC = helper.ParseUint32(records[11], 0)
				particleRender.Entries = append(particleRender.Entries, entry)
			}
			model.ParticleRenders = append(model.ParticleRenders, particleRender)
		}

		if mf.Name() == "particle_point.txt" {
			particlePoint := common.NewParticlePoint("")

			lines, err := helper.ReadFile(fmt.Sprintf("%s/%s", modelPath, mf.Name()))
			if err != nil {
				return fmt.Errorf("read model %s particle_point.txt: %w", model.Header.Name, err)
			}
			for i, line := range lines {
				if i == 0 {
					continue
				}
				if len(line) == 0 {
					continue
				}
				records := strings.Split(line, "|")
				if len(records) != 5 {
					return fmt.Errorf("particle_point.txt line %d: expected 5 records, got %d", i, len(records))
				}

				if records[0] == "id" {
					particlePoint.Header.Name = records[1]
					continue
				}

				entry := common.ParticlePointEntry{}
				entry.Name = records[0]
				entry.Bone = records[1]
				vec3 := strings.Split(records[2], ",")
				entry.Translation.X = helper.ParseFloat32(vec3[0], 0)
				entry.Translation.Y = helper.ParseFloat32(vec3[1], 0)
				entry.Translation.Z = helper.ParseFloat32(vec3[2], 0)
				vec3 = strings.Split(records[3], ",")
				entry.Rotation.X = helper.ParseFloat32(vec3[0], 0)
				entry.Rotation.Y = helper.ParseFloat32(vec3[1], 0)
				entry.Rotation.Z = helper.ParseFloat32(vec3[2], 0)
				vec3 = strings.Split(records[4], ",")
				entry.Scale.X = helper.ParseFloat32(vec3[0], 0)
				entry.Scale.Y = helper.ParseFloat32(vec3[1], 0)
				entry.Scale.Z = helper.ParseFloat32(vec3[2], 0)

				particlePoint.Entries = append(particlePoint.Entries, entry)
			}
			model.ParticlePoints = append(model.ParticlePoints, particlePoint)
		}
	}
	quail.Models = append(quail.Models, model)
	return nil
}

func (quail *Quail) dirParseMaterial(path string, name string) error {
	material := &common.Material{
		ShaderName: "Opaque_MaxCB1.fx",
	}
	material.Name = strings.TrimSuffix(name, ".material")
	_, ok := quail.materialCache[material.Name]
	if ok {
		// ignore materials already loaded
		return nil
	}
	lines, err := helper.ReadFile(fmt.Sprintf("%s/%s/property.txt", path, name))
	if err != nil {
		return fmt.Errorf("read material %s: %w", material.Name, err)
	}

	for i, line := range lines {
		if i == 0 {
			continue
		}
		if len(line) == 0 {
			continue
		}
		records := strings.Split(line, "|")
		recordType := strings.ToLower(records[0])
		if recordType == "shadername" {
			if len(records[1]) < 2 {
				continue
			}
			if records[1] == "None" {
				continue
			}
			material.ShaderName = records[1]
			continue
		}
		if len(records) != 3 {
			return fmt.Errorf("material %s line %d: expected 3 records, got %d", material.Name, i, len(records))
		}

		property := &common.MaterialProperty{}
		property.Name = records[0]
		property.Value = records[1]
		property.Category = helper.ParseUint32(records[2], 0)
		if property.Category == 2 && strings.Contains(strings.ToLower(property.Name), "texture") {
			property.Data, err = os.ReadFile(fmt.Sprintf("%s/%s/%s", path, name, property.Value))
			if err != nil {
				return fmt.Errorf("read material %s texture %s: %w", material.Name, property.Value, err)
			}
		}

		material.Properties = append(material.Properties, property)
	}

	quail.materialCache[material.Name] = material
	return nil
}

func (quail *Quail) dirParseAni(path string, name string) error {
	ani := &common.Animation{}
	ani.Header.Name = strings.TrimSuffix(name, ".ani")
	aniPath := fmt.Sprintf("%s/%s", path, name)
	aniFiles, err := os.ReadDir(aniPath)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", aniPath, err)
	}

	for _, af := range aniFiles {
		if af.Name() == "animation.txt" {
			lines, err := helper.ReadFile(fmt.Sprintf("%s/%s", aniPath, af.Name()))
			if err != nil {
				return fmt.Errorf("read ani %s animation.txt: %w", ani.Header.Name, err)
			}
			for i, line := range lines {
				if i == 0 {
					continue
				}
				if len(line) == 0 {
					continue
				}
				records := strings.Split(line, "|")
				if len(records) != 2 {
					return fmt.Errorf("animation.txt line %d: expected 2 records, got %d", i, len(records))
				}
				if records[0] == "is_strict" {
					ani.IsStrict = helper.ParseBool(records[1], false)
					continue
				}
			}
			continue
		}

		bone := &common.BoneAnimation{}
		bone.Name = af.Name()
		if strings.Contains(bone.Name, ".") {
			bone.Name = strings.Split(bone.Name, ".")[0]
		}
		lines, err := helper.ReadFile(fmt.Sprintf("%s/%s", aniPath, af.Name()))
		if err != nil {
			return fmt.Errorf("read ani %s %s: %w", ani.Header.Name, af.Name(), err)
		}
		for i, line := range lines {
			if i == 0 {
				continue
			}
			if len(line) == 0 {
				continue
			}
			records := strings.Split(line, "|")
			//milliseconds|rotation|scale|translation
			if len(records) != 4 {
				return fmt.Errorf("%s line %d: expected 4 records, got %d", af.Name(), i, len(records))
			}
			frame := &common.BoneAnimationFrame{}
			frame.Milliseconds = helper.ParseUint32(records[0], 0)
			vec4 := strings.Split(records[1], ",")
			frame.Rotation.X = helper.ParseFloat32(vec4[0], 0)
			frame.Rotation.Y = helper.ParseFloat32(vec4[1], 0)
			frame.Rotation.Z = helper.ParseFloat32(vec4[2], 0)
			frame.Rotation.W = helper.ParseFloat32(vec4[3], 0)
			vec3 := strings.Split(records[2], ",")
			frame.Scale.X = helper.ParseFloat32(vec3[0], 0)
			frame.Scale.Y = helper.ParseFloat32(vec3[1], 0)
			frame.Scale.Z = helper.ParseFloat32(vec3[2], 0)
			vec3 = strings.Split(records[3], ",")
			frame.Translation.X = helper.ParseFloat32(vec3[0], 0)
			frame.Translation.Y = helper.ParseFloat32(vec3[1], 0)
			frame.Translation.Z = helper.ParseFloat32(vec3[2], 0)
			bone.Frames = append(bone.Frames, frame)
			bone.FrameCount++
		}

		ani.Bones = append(ani.Bones, bone)
		continue

	}

	return nil
}
