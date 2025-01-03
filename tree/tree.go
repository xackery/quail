package tree

import (
	"fmt"
	"io"

	"github.com/xackery/quail/raw"
)

// Dump dumps a tree to a writer
func Dump(src interface{}, w io.Writer) error {
	var err error
	switch val := src.(type) {
	case *raw.Wld:
		err = wldDump(val, w)
		if err != nil {
			return fmt.Errorf("wld dump: %w", err)
		}
	case *raw.Bmp:
		return nil
	default:
		return fmt.Errorf("unknown type %T", val)
	}
	return nil
}
