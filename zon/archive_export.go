package zon

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/mds"
	"github.com/xackery/quail/mod"
	"github.com/xackery/quail/ter"
)

// ArchiveExport writes contents to outArchive
func (e *ZON) ArchiveExport(outArchive common.ArchiveReadWriter) error {
	if outArchive == nil {
		return fmt.Errorf("no archive loaded")
	}

	var err error

	for _, model := range e.models {
		modelData, err := e.archive.File(model.name)
		if err != nil {
			return fmt.Errorf("model file %s: %w", model.name, err)
		}

		switch filepath.Ext(model.name) {
		case ".ter":
			baseName := strings.TrimPrefix(helper.BaseName(model.name), "ter_")
			e, err := ter.New(baseName, e.archive)
			if err != nil {
				return fmt.Errorf("ter.NewEQG: %w", err)
			}
			err = e.Load(bytes.NewReader(modelData))
			if err != nil {
				return fmt.Errorf("ter load %s: %w", baseName, err)
			}

			err = e.ArchiveExport(outArchive)
			if err != nil {
				return fmt.Errorf("ter archive export %s: %w", baseName, err)
			}
		case ".mod":
			baseName := strings.TrimPrefix(helper.BaseName(model.name), "ter_")
			e, err := mod.New(baseName, e.archive)
			if err != nil {
				return fmt.Errorf("mod new: %w", err)
			}
			err = e.Load(bytes.NewReader(modelData))
			if err != nil {
				continue
				//return fmt.Errorf("mod load %s: %w", baseName, err)
			}
			err = e.ArchiveExport(outArchive)
			if err != nil {
				return fmt.Errorf("mod archive export %s: %w", baseName, err)
			}
		case ".mds":
			baseName := strings.TrimPrefix(helper.BaseName(model.name), "ter_")
			e, err := mds.New(baseName, e.archive)
			if err != nil {
				return fmt.Errorf("mds new: %w", err)
			}
			err = e.Load(bytes.NewReader(modelData))
			if err != nil {
				return fmt.Errorf("mds load %s: %w", baseName, err)
			}
			err = e.ArchiveExport(outArchive)
			if err != nil {
				return fmt.Errorf("mds archive export %s: %w", baseName, err)
			}

		default:
			return fmt.Errorf("unsupported model: %s", model.name)
		}

	}

	buf := bytes.NewBuffer(nil)
	err = e.Save(buf)
	if err != nil {
		return fmt.Errorf("zon save: %w", err)
	}

	err = outArchive.WriteFile(e.name+".zon", buf.Bytes())
	if err != nil {
		return fmt.Errorf("write %s.zon: %w", e.name, err)
	}

	return nil
}
