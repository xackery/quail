package quail

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

// EQGImport imports the quail target to an EQG file
func (e *Quail) EQGImport(path string) error {
	pfs, err := pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("eqg load: %w", err)
	}
	defer pfs.Close()

	particlePoints := []*common.ParticlePoint{}
	particleRenders := []*common.ParticleRender{}

	lits := []*common.RGBA{}
	for _, file := range pfs.Files() {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		reader := raw.New(ext)
		if reader == nil {
			return fmt.Errorf("unknown extension %s", ext)
		}
		err = reader.Read(bytes.NewReader(file.Data()))
		if err != nil {
			return fmt.Errorf("read %s: %w", file.Name(), err)
		}
	}

	if e.Zone != nil {
		e.Zone.Lits = lits
	}

	for _, point := range particlePoints {
		isFound := false
		for _, model := range e.Models {
			if strings.EqualFold(model.Header.Name, point.Header.Name) {
				isFound = true
				model.ParticlePoints = append(model.ParticlePoints, point)
				break
			}
		}
		if !isFound {
			log.Warnf("particle point %s not found in model", point.Header.Name)
		}
	}

	for _, render := range particleRenders {
		isFound := false
		for _, model := range e.Models {
			if strings.EqualFold(model.Header.Name, render.Header.Name) {
				isFound = true
				model.ParticleRenders = append(model.ParticleRenders, render)
				break
			}
		}
		if !isFound {
			log.Warnf("particle render %s not found in model", render.Header.Name)
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
