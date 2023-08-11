package quail

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergeymakinen/go-bmp"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/metadata/ani"
	"github.com/xackery/quail/model/metadata/prt"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/zon"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/pfs/s3d"
	"github.com/xackery/quail/quail/def"
	"github.com/xackery/quail/tag"
)

// Import imports the quail target
func (e *Quail) PFSImport(path string) error {
	ext := filepath.Ext(path)

	switch ext {
	case ".eqg":
		return e.EQGImport(path)
	case ".s3d":
		return e.S3DImport(path)
	default:
		return fmt.Errorf("unknown pfs type %s, valid options are eqg and pfs", ext)
	}
}

// EQGImport imports the quail target to an EQG file
func (e *Quail) EQGImport(path string) error {
	pfs, err := eqg.NewFile(path)
	if err != nil {
		return fmt.Errorf("eqg load: %w", err)
	}
	defer pfs.Close()

	particlePoints := []*def.ParticlePoint{}
	particleRenders := []*def.ParticleRender{}

	for _, file := range pfs.Files() {
		switch filepath.Ext(file.Name()) {
		case ".zon":
			e.Zone = &def.Zone{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".ZON"),
			}
			err = zon.Decode(e.Zone, bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodeZon %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.zon", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.zon.tags", "testdata", file.Name()))

		case ".pts":
			point := &def.ParticlePoint{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".PTS"),
			}
			err = pts.Decode(point, bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodePts %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.pts", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.pts.tags", "testdata", file.Name()))
			particlePoints = append(particlePoints, point)
		case ".prt":
			render := &def.ParticleRender{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".PRT"),
			}
			err = prt.Decode(render, bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodePrt %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.prt", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.prt.tags", "testdata", file.Name()))
			particleRenders = append(particleRenders, render)
		case ".ani":
			anim := &def.Animation{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".ANI"),
			}
			err = ani.Decode(anim, bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodeAni %s: %w", file.Name(), err)
			}

			os.WriteFile(fmt.Sprintf("%s/%s-raw.ani", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.ani.tags", "testdata", file.Name()))
			e.Animations = append(e.Animations, anim)
		case ".mod":
			mesh := &def.Mesh{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".MOD"),
			}
			err = mod.Decode(mesh, bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodeMod %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.mod", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.mod.tags", "testdata", file.Name()))

			e.Meshes = append(e.Meshes, mesh)
		case ".ter":
			mesh := &def.Mesh{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".MOD"),
			}
			err = TERDecode(mesh, bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodeter %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.ter", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.ter.tags", "testdata", file.Name()))

			e.Meshes = append(e.Meshes, mesh)
		case ".mds":
			mesh := &def.Mesh{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".MDS"),
			}
			err = mds.Decode(mesh, bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodeMds %s: %w", file.Name(), err)
			}
			e.Meshes = append(e.Meshes, mesh)
		}
	}

	for _, point := range particlePoints {
		isFound := false
		for _, mesh := range e.Meshes {
			if strings.EqualFold(mesh.Name, point.Name) {
				isFound = true
				mesh.ParticlePoints = append(mesh.ParticlePoints, point)
				break
			}
		}
		if !isFound {
			log.Warnf("particle point %s not found in mesh", point.Name)
		}
	}

	for _, render := range particleRenders {
		isFound := false
		for _, mesh := range e.Meshes {
			if strings.EqualFold(mesh.Name, render.Name) {
				isFound = true
				mesh.ParticleRenders = append(mesh.ParticleRenders, render)
				break
			}
		}
		if !isFound {
			log.Warnf("particle render %s not found in mesh", render.Name)
		}
	}

	materialCount := 0
	textureCount := 0
	for _, mesh := range e.Meshes {
		for _, material := range mesh.Materials {
			materialCount++
			for _, property := range material.Properties {
				if property.Category != 2 {
					continue
				}
				if !strings.Contains(strings.ToLower(property.Name), "texture") {
					continue
				}
				for _, file := range pfs.Files() {
					if strings.EqualFold(file.Name(), property.Value) {
						property.Data = file.Data()
						textureCount++
					}
				}
			}
		}
	}

	log.Debugf("%s (eqg) loaded %d meshes, %d materials, %d texture files", filepath.Base(path), len(e.Meshes), materialCount, textureCount)
	return nil
}

// S3DImport imports the quail target to an S3D file
func (e *Quail) S3DImport(path string) error {
	pfs, err := s3d.NewFile(path)
	if err != nil {
		return fmt.Errorf("s3d load: %w", err)
	}
	defer pfs.Close()

	for _, file := range pfs.Files() {
		switch filepath.Ext(file.Name()) {
		case ".wld":
			if !strings.HasSuffix(file.Name(), ".wld") {
				continue
			}

			log.Debugf("testing %s", file.Name())

			meshes, err := WLDDecode(bytes.NewReader(file.Data()), pfs)
			if err != nil {
				return fmt.Errorf("wldDecode %s: %w", file.Name(), err)
			}
			e.Meshes = append(e.Meshes, meshes...)
		}
	}

	if len(e.Meshes) == 1 && e.Meshes[0].Name == "" {
		e.Meshes[0].Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	materialCount := 0
	textureCount := 0
	for _, mesh := range e.Meshes {
		log.Debugf("mesh %s has %d materials", mesh.Name, len(mesh.Materials))
		for _, material := range mesh.Materials {
			materialCount++
			for _, property := range material.Properties {
				if property.Category != 2 {
					continue
				}
				if !strings.Contains(strings.ToLower(property.Name), "texture") {
					continue
				}
				for _, file := range pfs.Files() {
					if !strings.EqualFold(file.Name(), property.Value) {
						continue
					}
					property.Data = file.Data()
					textureCount++

					if string(property.Data[0:3]) == "DDS" {
						property.Value = strings.TrimSuffix(property.Value, filepath.Ext(property.Value)) + ".dds"
						material.Name = strings.TrimSuffix(strings.TrimSuffix(material.Name, ".BMP"), ".bmp")
						continue
					}

					if filepath.Ext(strings.ToLower(property.Value)) != ".bmp" {
						continue
					}
					img, err := bmp.Decode(bytes.NewReader(file.Data()))
					if err != nil {
						return fmt.Errorf("bmp decode: %w", err)
					}
					buf := new(bytes.Buffer)
					// convert to png
					err = png.Encode(buf, img)
					if err != nil {
						return fmt.Errorf("png encode: %w", err)
					}
					property.Value = strings.TrimSuffix(property.Value, filepath.Ext(property.Value)) + ".png"
					material.Name = strings.TrimSuffix(material.Name, ".bmp")
				}

			}
		}
	}

	log.Debugf("%s (s3d) loaded %d meshes, %d materials, %d texture files", filepath.Base(path), len(e.Meshes), materialCount, textureCount)
	return nil
}
