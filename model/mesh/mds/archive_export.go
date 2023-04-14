package mds

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/pfs/archive"
)

// ArchiveExport writes contents to outArchive
func (e *MDS) ArchiveExport(outArchive archive.Writer) error {
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

	return nil
}
