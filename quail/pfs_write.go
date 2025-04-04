package quail

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/wce"
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
	archive, err := pfs.New(path)
	if err != nil {
		return fmt.Errorf("eqg new: %w", err)
	}
	defer archive.Close()

	if e.Wld == nil {
		return fmt.Errorf("no wld found")
	}

	err = e.Wld.WriteEqgRaw(archive)
	if err != nil {
		return fmt.Errorf("write eqg: %w", err)
	}

	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer w.Close()

	for fileName, assetData := range e.Assets {
		err := archive.Add(fileName, assetData)
		if err != nil {
			return fmt.Errorf("addAsset %s: %w", fileName, err)
		}
	}

	err = archive.Write(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}

	noExtBaseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	dirPath := filepath.Dir(path)

	sideFileOut := ""
	numSideFiles := 0
	ok, err := e.Wld.WriteSingleFile(wce.SideFileZon, filepath.Join(dirPath, noExtBaseName+".zon"))
	if err != nil {
		return fmt.Errorf("write side file .zon: %w", err)
	}
	if ok {
		sideFileOut += ".zon, "
		numSideFiles++
	}

	if len(sideFileOut) > 0 {
		sideFileOut = strings.TrimSuffix(sideFileOut, ", ")
		sideFileOut = fmt.Sprintf(" and side file%s: %s", helper.Pluralize(numSideFiles), sideFileOut)
	}

	fmt.Printf("Wrote %s with %d entries%s\n", path, archive.Len(), sideFileOut)

	return nil
}

// S3DExport exports the quail target to an S3D file
func (e *Quail) S3DExport(fileVersion uint32, pfsVersion int, path string) error {
	archive, err := pfs.New(path)
	if err != nil {
		return fmt.Errorf("eqg new: %w", err)
	}
	defer archive.Close()

	isSomethingWritten := false
	if e.Wld != nil {
		buf := &bytes.Buffer{}

		err := e.Wld.WriteWldRaw(buf)
		if err != nil {
			return fmt.Errorf("write s3d: %w", err)
		}

		err = archive.Add(e.Wld.FileName, buf.Bytes())
		if err != nil {
			return fmt.Errorf("addWld %s: %w", e.Wld.FileName, err)
		}
		isSomethingWritten = true
	}

	if e.WldObject != nil {
		buf := &bytes.Buffer{}

		err := e.WldObject.WriteWldRaw(buf)
		if err != nil {
			return fmt.Errorf("write s3d object: %w", err)
		}

		err = archive.Add("objects.wld", buf.Bytes())
		if err != nil {
			return fmt.Errorf("addWld %s: %w", e.Wld.FileName, err)
		}
		isSomethingWritten = true
	}

	if e.WldLights != nil {
		buf := &bytes.Buffer{}

		err := e.WldLights.WriteWldRaw(buf)
		if err != nil {
			return fmt.Errorf("write s3d lights: %w", err)
		}

		err = archive.Add("lights.wld", buf.Bytes())
		if err != nil {
			return fmt.Errorf("addWld %s: %w", e.Wld.FileName, err)
		}
		isSomethingWritten = true
	}

	for fileName, assetData := range e.Assets {
		err := archive.Add(fileName, assetData)
		if err != nil {
			return fmt.Errorf("addAsset %s: %w", fileName, err)
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

	err = archive.Write(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	fmt.Printf("Wrote %s with %d entries\n", path, archive.Len())

	return nil
}
