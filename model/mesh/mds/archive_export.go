package mds

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xackery/quail/common"
)

// ArchiveExport writes contents to outArchive
func (e *MDS) ArchiveExport(outArchive common.ArchiveWriter) error {
	if outArchive == nil {
		return fmt.Errorf("no archive loaded")
	}

	buf := bytes.NewBuffer(nil)
	err := e.Encode(buf)
	if err != nil {
		return fmt.Errorf("mds encode: %w", err)
	}

	err = outArchive.WriteFile(e.name+".mds", buf.Bytes())
	if err != nil {
		return fmt.Errorf("write %s.mds: %w", e.name, err)
	}

	for _, material := range e.materials {
		for _, p := range material.Properties {
			if p.Category != 2 {
				continue
			}
			if !strings.EqualFold(p.Name, "e_texturediffuse0") &&
				!strings.EqualFold(p.Name, "e_texturenormal0") {
				continue
			}

			data, err := e.archive.File(p.Value)
			if err != nil {
				return fmt.Errorf("file: %w", err)
			}
			err = outArchive.WriteFile(p.Value, data)
			if err != nil {
				return fmt.Errorf("writeFile: %w", err)
			}
		}
	}

	return nil
}
