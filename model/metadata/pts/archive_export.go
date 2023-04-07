package pts

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/pfs/archive"
)

// ArchiveExport writes contents to outArchive
func (e *PTS) ArchiveExport(outArchive archive.Writer) error {
	if outArchive == nil {
		return fmt.Errorf("no archive loaded")
	}

	buf := bytes.NewBuffer(nil)
	err := e.Encode(buf)
	if err != nil {
		return fmt.Errorf("pts encode: %w", err)
	}

	err = outArchive.WriteFile(e.name+".pts", buf.Bytes())
	if err != nil {
		return fmt.Errorf("write %s.pts: %w", e.name, err)
	}

	return nil
}
