package mod

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail/common"
)

// ArchiveExport writes contents to outArchive
func (e *MOD) ArchiveExport(outArchive common.ArchiveWriter) error {
	if outArchive == nil {
		return fmt.Errorf("no archive loaded")
	}

	buf := bytes.NewBuffer(nil)
	err := e.Save(buf)
	if err != nil {
		return fmt.Errorf("mod save: %w", err)
	}

	err = outArchive.WriteFile(e.name+".mod", buf.Bytes())
	if err != nil {
		return fmt.Errorf("write %s.mod: %w", e.name, err)
	}

	return nil
}
