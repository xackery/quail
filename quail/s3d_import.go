package quail

import (
	"bytes"
	"fmt"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/malashin/dds"
	"github.com/sergeymakinen/go-bmp"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

// S3DImport imports the quail target to an S3D file
func (e *Quail) S3DImport(path string) error {
	pfs, err := pfs.NewFile(path)
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

			wld := &raw.Wld{}
			err = wld.Read(bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("wldDecode %s: %w", file.Name(), err)
			}
			err = e.RawRead(wld)
			if err != nil {
				return fmt.Errorf("rawRead %s: %w", file.Name(), err)
			}
		}
	}

	if len(e.Models) == 1 && e.Models[0].Header.Name == "" {
		e.Models[0].Header.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	materialCount := 0
	textureCount := 0
	for _, model := range e.Models {
		log.Debugf("model %s has %d materials", model.Header.Name, len(model.Materials))
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
						// change to png, blender doesn't like EQ dds
						img, err := dds.Decode(bytes.NewReader(property.Data))
						if err != nil {
							return fmt.Errorf("dds decode: %w", err)
						}
						buf := new(bytes.Buffer)
						err = png.Encode(buf, img)
						if err != nil {
							return fmt.Errorf("png encode: %w", err)
						}
						if strings.HasSuffix(strings.ToLower(material.Name), ".bmp") {
							material.Name = strings.TrimSuffix(material.Name, ".bmp")
						}
						property.Data = buf.Bytes()
						property.Value = strings.TrimSuffix(property.Value, filepath.Ext(property.Value)) + ".png"
						continue
					}

					if filepath.Ext(strings.ToLower(property.Value)) != ".bmp" {
						continue
					}
					img, err := bmp.Decode(bytes.NewReader(file.Data()))
					if err != nil {
						return fmt.Errorf("bmp read: %w", err)
					}
					buf := new(bytes.Buffer)
					// convert to png
					err = png.Encode(buf, img)
					if err != nil {
						return fmt.Errorf("png encode: %w", err)
					}
					property.Value = strings.TrimSuffix(property.Value, filepath.Ext(property.Value)) + ".png"
					material.Name = strings.TrimSuffix(material.Name, ".bmp")
					property.Data = buf.Bytes()
				}

			}
		}
	}

	log.Debugf("%s (s3d) loaded %d models, %d materials, %d texture files", filepath.Base(path), len(e.Models), materialCount, textureCount)
	return nil
}
