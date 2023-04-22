package wld

import (
	"fmt"
	"io"
)

// 0x00 0
func (e *WLD) defaultRead(r io.ReadSeeker, fragmentOffset int) error {
	return fmt.Errorf("default fallback, unsupported fragment")
}
