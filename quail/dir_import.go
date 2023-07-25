package quail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/quail/def"
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
		case ".mesh":
			mesh := &def.Mesh{
				FileType: "mod", // default to mod
			}
			mesh.Name = strings.TrimSuffix(qf.Name(), ".mesh")
			meshPath := fmt.Sprintf("%s/%s", path, qf.Name())
			meshFiles, err := os.ReadDir(meshPath)
			if err != nil {
				return fmt.Errorf("read dir %s: %w", meshPath, err)
			}
			for _, mf := range meshFiles {
				if mf.Name() == "triangle.txt" {
					triangleData, err := os.ReadFile(fmt.Sprintf("%s/%s", meshPath, mf.Name()))
					if err != nil {
						return fmt.Errorf("read mesh %s triangle.txt: %w", mesh.Name, err)
					}
					lines := strings.Split(string(triangleData), "\n")
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
						triangle := def.Triangle{}
						vec3 := strings.Split(records[0], ",")
						triangle.Index.X = helper.ParseUint32(vec3[0], 0)
						triangle.Index.Y = helper.ParseUint32(vec3[1], 0)
						triangle.Index.Z = helper.ParseUint32(vec3[2], 0)
						triangle.Flag = helper.ParseUint32(records[1], 0)
						triangle.MaterialName = records[2]
						mesh.Triangles = append(mesh.Triangles, triangle)
					}
				}
				if mf.Name() == "vertex.txt" {
					vertexData, err := os.ReadFile(fmt.Sprintf("%s/%s", meshPath, mf.Name()))
					if err != nil {
						return fmt.Errorf("read mesh %s vertex.txt: %w", mesh.Name, err)
					}
					lines := strings.Split(string(vertexData), "\n")
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
						vertex := def.Vertex{}
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

						mesh.Vertices = append(mesh.Vertices, vertex)
					}
				}

				if mf.Name() == "bone.txt" {
					boneData, err := os.ReadFile(fmt.Sprintf("%s/%s", meshPath, mf.Name()))
					if err != nil {
						return fmt.Errorf("read mesh %s bone.txt: %w", mesh.Name, err)
					}
					lines := strings.Split(string(boneData), "\n")
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
						bone := def.Bone{}
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

						mesh.Bones = append(mesh.Bones, bone)
					}
				}

				if mf.Name() == "particle_render.txt" {
					particleRender := &def.ParticleRender{}

					particleRenderData, err := os.ReadFile(fmt.Sprintf("%s/%s", meshPath, mf.Name()))
					if err != nil {
						return fmt.Errorf("read mesh %s particle_render.txt: %w", mesh.Name, err)
					}
					lines := strings.Split(string(particleRenderData), "\n")
					for i, line := range lines {
						if i == 0 {
							continue
						}
						if len(line) == 0 {
							continue
						}
						records := strings.Split(line, "|")
						if len(records) != 8 {
							return fmt.Errorf("particle_render.txt line %d: expected 8 records, got %d", i, len(records))
						}

						entry := &def.ParticleRenderEntry{}
						entry.ID = helper.ParseUint32(records[0], 0)
						entry.ID2 = helper.ParseUint32(records[1], 0)
						entry.ParticlePoint = records[2]
						b5 := strings.Split(records[3], ",")
						entry.UnknownA[0] = helper.ParseUint32(b5[0], 0)
						entry.UnknownA[1] = helper.ParseUint32(b5[1], 0)
						entry.UnknownA[2] = helper.ParseUint32(b5[2], 0)
						entry.UnknownA[3] = helper.ParseUint32(b5[3], 0)
						entry.UnknownA[4] = helper.ParseUint32(b5[4], 0)

						entry.Duration = helper.ParseUint32(records[4], 0)
						entry.UnknownB = helper.ParseUint32(records[5], 0)
						entry.UnknownFFFFFFFF = helper.ParseInt32(records[6], 0)
						entry.UnknownC = helper.ParseUint32(records[7], 0)
						particleRender.Entries = append(particleRender.Entries, entry)
					}
					mesh.ParticleRenders = append(mesh.ParticleRenders, particleRender)
				}

				if mf.Name() == "particle_point.txt" {
					particlePoint := &def.ParticlePoint{}

					particlePointData, err := os.ReadFile(fmt.Sprintf("%s/%s", meshPath, mf.Name()))
					if err != nil {
						return fmt.Errorf("read mesh %s particle_point.txt: %w", mesh.Name, err)
					}
					lines := strings.Split(string(particlePointData), "\n")
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
							particlePoint.Name = records[1]
							continue
						}

						entry := def.ParticlePointEntry{}
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
					mesh.ParticlePoints = append(mesh.ParticlePoints, particlePoint)
				}
				if mf.IsDir() && strings.HasSuffix(mf.Name(), ".material") {
					material := &def.Material{
						ShaderName: "Opaque_MaxCB1.fx",
					}
					material.Name = strings.TrimSuffix(mf.Name(), ".material")
					materialData, err := os.ReadFile(fmt.Sprintf("%s/%s/property.txt", meshPath, mf.Name()))
					if err != nil {
						return fmt.Errorf("read mesh %s material %s: %w", mesh.Name, mf.Name(), err)
					}
					lines := strings.Split(string(materialData), "\n")
					for i, line := range lines {
						if i == 0 {
							continue
						}
						if len(line) == 0 {
							continue
						}
						records := strings.Split(line, "|")
						if records[0] == "shaderName" {
							material.ShaderName = records[1]
							continue
						}
						if len(records) != 3 {
							return fmt.Errorf("material %s line %d: expected 3 records, got %d", mf.Name(), i, len(records))
						}

						property := &def.MaterialProperty{}
						property.Name = records[0]
						property.Value = records[1]
						property.Category = helper.ParseUint32(records[2], 0)
						if property.Category == 2 && strings.Contains(strings.ToLower(property.Name), "texture") {
							property.Data, err = os.ReadFile(fmt.Sprintf("%s/%s/%s", meshPath, mf.Name(), property.Value))
							if err != nil {
								return fmt.Errorf("read mesh %s material %s texture %s: %w", mesh.Name, mf.Name(), property.Value, err)
							}
						}

						material.Properties = append(material.Properties, property)
					}
					mesh.Materials = append(mesh.Materials, material)
				}
			}
			quail.Meshes = append(quail.Meshes, mesh)
		}
	}
	return nil
}
