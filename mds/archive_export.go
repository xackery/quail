package mds

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
)

// ArchiveExport writes contents to outArchive
func (e *MDS) ArchiveExport(outArchive common.ArchiveWriter) error {
	if outArchive == nil {
		return fmt.Errorf("no archive loaded")
	}

	buf := bytes.NewBuffer(nil)
	err := e.Save(buf)
	if err != nil {
		return fmt.Errorf("mds save: %w", err)
	}

	err = outArchive.WriteFile(e.name+".mds", buf.Bytes())
	if err != nil {
		return fmt.Errorf("write %s.mds: %w", e.name, err)
	}

	return nil
}
