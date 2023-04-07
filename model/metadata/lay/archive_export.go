package lay

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/pfs/archive"
)

// ArchiveExport writes contents to outArchive
func (e *LAY) ArchiveExport(outArchive archive.Writer) error {
	if outArchive == nil {
		return fmt.Errorf("no archive loaded")
	}

	buf := bytes.NewBuffer(nil)
	err := e.Encode(buf)
	if err != nil {
		return fmt.Errorf("lay encode: %w", err)
	}

	err = outArchive.WriteFile(e.name+".lay", buf.Bytes())
	if err != nil {
		return fmt.Errorf("write %s.lay: %w", e.name, err)
	}

	return nil
}
