package quail

import (
	"bytes"
	"fmt"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/sergeymakinen/go-bmp"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs/s3d"
)

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

			modeles, err := WLDDecode(bytes.NewReader(file.Data()), pfs)
			if err != nil {
				return fmt.Errorf("wldDecode %s: %w", file.Name(), err)
			}
			e.Models = append(e.Models, modeles...)
		}
	}

	if len(e.Models) == 1 && e.Models[0].Name == "" {
		e.Models[0].Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	materialCount := 0
	textureCount := 0
	for _, model := range e.Models {
		log.Debugf("model %s has %d materials", model.Name, len(model.Materials))
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

	log.Debugf("%s (s3d) loaded %d modeles, %d materials, %d texture files", filepath.Base(path), len(e.Models), materialCount, textureCount)
	return nil
}
