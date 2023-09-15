package quail

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/mesh/ter"
	"github.com/xackery/quail/model/metadata/ani"
	"github.com/xackery/quail/model/metadata/lay"
	"github.com/xackery/quail/model/metadata/prt"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/zon"
	"github.com/xackery/quail/pfs/eqg"
	"github.com/xackery/quail/tag"
)

// EQGImport imports the quail target to an EQG file
func (e *Quail) EQGImport(path string) error {
	pfs, err := eqg.NewFile(path)
	if err != nil {
		return fmt.Errorf("eqg load: %w", err)
	}
	defer pfs.Close()

	particlePoints := []*common.ParticlePoint{}
	particleRenders := []*common.ParticleRender{}

	lits := []*common.RGBA{}
	for _, file := range pfs.Files() {
		switch filepath.Ext(file.Name()) {
		case ".zon":
			e.Zone = &common.Zone{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".ZON"),
			}
			err = zon.Decode(e.Zone, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), e.Zone.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodeZon %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.zon", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.zon.tags", "testdata", file.Name()))
		case ".pts":
			point := &common.ParticlePoint{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".PTS"),
			}
			err = pts.Decode(point, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), point.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodePts %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.pts", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.pts.tags", "testdata", file.Name()))
			particlePoints = append(particlePoints, point)
		case ".prt":
			render := &common.ParticleRender{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".PRT"),
			}
			err = prt.Decode(render, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), render.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodePrt %s: %w", file.Name(), err)
			}

			os.WriteFile(fmt.Sprintf("%s/%s-raw.prt", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.prt.tags", "testdata", file.Name()))
			particleRenders = append(particleRenders, render)
		case ".lay":
			model := &common.Model{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".LAY"),
			}
			err = lay.Decode(model, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), model.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodePrt %s: %w", file.Name(), err)
			}
		case ".lit":
			//err = lit.Decode(lits, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), 1, file.Name(), filepath.Base(path))
			}
			//if err != nil {
			//	return fmt.Errorf("decodeLit %s: %w", file.Name(), err)
			//}
		case ".ani":
			anim := &common.Animation{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".ANI"),
			}
			err = ani.Decode(anim, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), anim.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodeAni %s: %w", file.Name(), err)
			}

			os.WriteFile(fmt.Sprintf("%s/%s-raw.ani", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.ani.tags", "testdata", file.Name()))
			e.Animations = append(e.Animations, anim)
		case ".mod":
			model := &common.Model{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".MOD"),
			}
			err = mod.Decode(model, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), model.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodeMod %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.mod", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.mod.tags", "testdata", file.Name()))

			e.Models = append(e.Models, model)
		case ".ter":
			model := &common.Model{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".MOD"),
			}
			err = ter.Decode(model, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), model.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodeter %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.ter", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.ter.tags", "testdata", file.Name()))

			e.Models = append(e.Models, model)
		case ".mds":
			model := &common.Model{
				Name: strings.TrimSuffix(strings.ToUpper(file.Name()), ".MDS"),
			}
			err = mds.Decode(model, bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), model.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodeMds %s: %w", file.Name(), err)
			}
			e.Models = append(e.Models, model)
		case ".dds":
		case ".bmp":
		case ".png":
		default:
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|?|%s|%s\n", filepath.Ext(file.Name()), file.Name(), filepath.Base(path))
			}
		}

	}

	if e.Zone != nil {
		e.Zone.Lits = lits
	}

	for _, point := range particlePoints {
		isFound := false
		for _, model := range e.Models {
			if strings.EqualFold(model.Name, point.Name) {
				isFound = true
				model.ParticlePoints = append(model.ParticlePoints, point)
				break
			}
		}
		if !isFound {
			log.Warnf("particle point %s not found in model", point.Name)
		}
	}

	for _, render := range particleRenders {
		isFound := false
		for _, model := range e.Models {
			if strings.EqualFold(model.Name, render.Name) {
				isFound = true
				model.ParticleRenders = append(model.ParticleRenders, render)
				break
			}
		}
		if !isFound {
			log.Warnf("particle render %s not found in model", render.Name)
		}
	}

	materialCount := 0
	textureCount := 0
	for _, model := range e.Models {
		for _, material := range model.Materials {
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

	log.Debugf("%s (eqg) loaded %d models, %d materials, %d texture files", filepath.Base(path), len(e.Models), materialCount, textureCount)
	return nil
}
