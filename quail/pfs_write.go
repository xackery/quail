package quail

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

// Write exports the quail target
func (e *Quail) PfsWrite(fileVersion uint32, pfsVersion int, path string) error {
	if len(path) == 0 {
		return fmt.Errorf("path is empty")
	}
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".eqg":
		return e.EQGExport(fileVersion, pfsVersion, path)
	case ".s3d":
		return e.S3DExport(fileVersion, pfsVersion, path)
	default:
		if len(ext) < 2 {
			return fmt.Errorf("unknown pfs type %s, valid options are eqg and pfs", path)
		}

		return fmt.Errorf("unknown pfs type %s, valid options are eqg and pfs", ext[1:])
	}
}

// EQGExport exports the quail target to an EQG file
func (e *Quail) EQGExport(fileVersion uint32, pfsVersion int, path string) error {
	pfs, err := pfs.New(path)
	if err != nil {
		return fmt.Errorf("eqg new: %w", err)
	}
	defer pfs.Close()

	if e.Zone != nil {
		buf := &bytes.Buffer{}

		zon := &raw.Zon{}
		err = e.RawWrite(zon)
		if err != nil {
			return fmt.Errorf("write quail->zon: %w", err)
		}

		err = zon.Write(buf)
		if err != nil {
			return fmt.Errorf("write zon: %w", err)
		}
		//os.WriteFile(fmt.Sprintf("%s/%s-raw-out.zon", "testdata", e.Zone.Header.Name), buf.Bytes(), 0644)
		//tag.Write(fmt.Sprintf("%s/%s-raw-out.zon.tags", "testdata", e.Zone.Header.Name))

		err = pfs.Add(fmt.Sprintf("%s.zon", e.Zone.Header.Name), buf.Bytes())
		if err != nil {
			return fmt.Errorf("addZon %s: %w", e.Zone.Header.Name, err)
		}
	}

	for _, entry := range e.Models {
		buf := &bytes.Buffer{}
		switch entry.FileType {
		case "ter":
			ter := &raw.Ter{}
			err = e.RawWrite(ter)
			if err != nil {
				return fmt.Errorf("write quail->ter: %w", err)
			}
			err = ter.Write(buf)
			if err != nil {
				return fmt.Errorf("ter.Write %s: %w", entry.Header.Name, err)
			}
		default:
			return fmt.Errorf("unknown filetype %s", entry.FileType)
		}

		//os.WriteFile(fmt.Sprintf("%s/%s-raw-out.%s", "testdata", entry.Header.Name, entry.FileType), buf.Bytes(), 0644)
		//tag.Write(fmt.Sprintf("%s/%s-raw-out.%s.tags", "testdata", entry.Header.Name, entry.FileType))

		err = pfs.Add(fmt.Sprintf("%s.%s", entry.Header.Name, entry.FileType), buf.Bytes())
		if err != nil {
			return fmt.Errorf("addMod %s: %w", entry.Header.Name, err)
		}
		for _, material := range entry.Materials {
			for _, property := range material.Properties {
				if len(property.Data) == 0 {
					continue
				}
				err = pfs.Add(property.Value, property.Data)
				if err != nil {
					return fmt.Errorf("model %s addMaterial %s texture %s: %w", entry.Header.Name, material.Name, property.Value, err)
				}
			}
		}
	}

	for _, anim := range e.Animations {
		buf := &bytes.Buffer{}
		ani := &raw.Ani{}
		err = e.RawWrite(ani)
		if err != nil {
			return fmt.Errorf("write quail->ani: %w", err)
		}
		err = ani.Write(buf)
		if err != nil {
			return fmt.Errorf("encodeAni %s: %w", anim.Header.Name, err)
		}
		err = pfs.Add(fmt.Sprintf("%s.ani", anim.Header.Name), buf.Bytes())
		if err != nil {
			return fmt.Errorf("addMds %s: %w", anim.Header.Name, err)
		}
	}

	for _, model := range e.Models {
		for _, render := range model.ParticleRenders {
			buf := &bytes.Buffer{}
			prt := &raw.Prt{}
			err = e.RawWrite(prt)
			if err != nil {
				return fmt.Errorf("write quail->prt: %w", err)
			}
			err = prt.Write(buf)
			if err != nil {
				return fmt.Errorf("encodePrt %s: %w", render.Header.Name, err)
			}
			err = pfs.Add(fmt.Sprintf("%s.prt", render.Header.Name), buf.Bytes())
			if err != nil {
				return fmt.Errorf("addPtr %s: %w", render.Header.Name, err)
			}
		}

		for _, point := range model.ParticlePoints {
			buf := &bytes.Buffer{}
			pts := &raw.Pts{}
			err = e.RawWrite(pts)
			if err != nil {
				return fmt.Errorf("write quail->pts: %w", err)
			}
			err = pts.Write(buf)
			if err != nil {
				return fmt.Errorf("encodePts %s: %w", point.Header.Name, err)
			}
			err = pfs.Add(fmt.Sprintf("%s.pts", point.Header.Name), buf.Bytes())
			if err != nil {
				return fmt.Errorf("addPts %s: %w", point.Header.Name, err)
			}
		}
	}

	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer w.Close()

	err = pfs.Write(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	log.Debugf("wrote %s with %d entries", path, pfs.Len())
	return nil
}

// S3DExport exports the quail target to an S3D file
func (e *Quail) S3DExport(fileVersion uint32, pfsVersion int, path string) error {
	pfs, err := pfs.New(path)
	if err != nil {
		return fmt.Errorf("eqg new: %w", err)
	}
	defer pfs.Close()

	isSomethingWritten := false
	if e.wld != nil {
		buf := &bytes.Buffer{}

		err := e.wld.WriteRaw(buf)
		if err != nil {
			return fmt.Errorf("write wld: %w", err)
		}

		err = pfs.Add(e.wld.FileName, buf.Bytes())
		if err != nil {
			return fmt.Errorf("addWld %s: %w", e.wld.FileName, err)
		}
		isSomethingWritten = true
	}

	if e.wldObject != nil {
		buf := &bytes.Buffer{}

		err := e.wldObject.WriteRaw(buf)
		if err != nil {
			return fmt.Errorf("write wld: %w", err)
		}

		err = pfs.Add("objects.wld", buf.Bytes())
		if err != nil {
			return fmt.Errorf("addWld %s: %w", e.wld.FileName, err)
		}
		isSomethingWritten = true
	}

	if e.wldLights != nil {
		buf := &bytes.Buffer{}

		err := e.wldLights.WriteRaw(buf)
		if err != nil {
			return fmt.Errorf("write wld: %w", err)
		}

		err = pfs.Add("lights.wld", buf.Bytes())
		if err != nil {
			return fmt.Errorf("addWld %s: %w", e.wld.FileName, err)
		}
		isSomethingWritten = true
	}

	for fileName, textureData := range e.Textures {
		err := pfs.Add(fileName, textureData)
		if err != nil {
			return fmt.Errorf("addTexture %s: %w", fileName, err)
		}
		isSomethingWritten = true
	}

	if !isSomethingWritten {
		return fmt.Errorf("nothing written")
	}

	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer w.Close()

	err = pfs.Write(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	log.Debugf("wrote %s with %d entries", path, pfs.Len())

	return nil
}
