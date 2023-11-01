package quail

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/tag"
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
		switch filepath.Ext(file.Name()) {
		case ".zon":
			zon := &raw.Zon{}
			err = zon.Read(bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), e.Zone.Header.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodeZon %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.zon", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.zon.tags", "testdata", file.Name()))
		case ".pts":
			pts := &raw.Pts{}
			err = pts.Read(bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), pts.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodePts %s: %w", file.Name(), err)
			}
			os.WriteFile(fmt.Sprintf("%s/%s-raw.pts", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.pts.tags", "testdata", file.Name()))
			err = e.RawRead(pts)
			if err != nil {
				return fmt.Errorf("rawRead %s: %w", file.Name(), err)
			}
		case ".prt":
			prt := &raw.Prt{}
			err = prt.Read(bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), prt.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodePrt %s: %w", file.Name(), err)
			}

			os.WriteFile(fmt.Sprintf("%s/%s-raw.prt", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.prt.tags", "testdata", file.Name()))
			err = e.RawRead(prt)
			if err != nil {
				return fmt.Errorf("rawRead %s: %w", file.Name(), err)
			}
		case ".lay":
			model := common.NewModel(strings.TrimSuffix(strings.ToUpper(file.Name()), ".LAY"))
			lay := &raw.Lay{}
			err = lay.Read(bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodePrt %s: %w", file.Name(), err)
			}
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), model.Header.Version, file.Name(), filepath.Base(path))
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
			ani := &raw.Ani{}
			err = ani.Read(bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), ani.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodeAni %s: %w", file.Name(), err)
			}

			os.WriteFile(fmt.Sprintf("%s/%s-raw.ani", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.ani.tags", "testdata", file.Name()))
		case ".mod":
			mod := &raw.Mod{}
			err = mod.Read(bytes.NewReader(file.Data()))
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), mod.Version, file.Name(), filepath.Base(path))
			}
			if err != nil {
				return fmt.Errorf("decodeMod %s: %w", file.Name(), err)
			}

			os.WriteFile(fmt.Sprintf("%s/%s-raw.mod", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.mod.tags", "testdata", file.Name()))
		case ".ter":
			ter := &raw.Ter{}
			err = ter.Read(bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodeter %s: %w", file.Name(), err)
			}
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), ter.Version, file.Name(), filepath.Base(path))
			}

			os.WriteFile(fmt.Sprintf("%s/%s-raw.ter", "testdata", file.Name()), file.Data(), 0644)
			tag.Write(fmt.Sprintf("%s/%s-raw.ter.tags", "testdata", file.Name()))
		case ".mds":
			mds := &raw.Mds{}
			err = mds.Read(bytes.NewReader(file.Data()))
			if err != nil {
				return fmt.Errorf("decodeMds %s: %w", file.Name(), err)
			}
			if e.IsExtensionVersionDump {
				fmt.Printf("%s|%d|%s|%s\n", filepath.Ext(file.Name()), mds.Version, file.Name(), filepath.Base(path))
			}
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
