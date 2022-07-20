package lit

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/quail/dump"
)

func (e *LIT) Decode(r io.ReadSeeker) error {
	var err error

	lightCount := uint32(0)
	err = binary.Read(r, binary.LittleEndian, &lightCount)
	if err != nil {
		return fmt.Errorf("read lightCount: %w", err)
	}
	dump.Hex(lightCount, "lightCount=%d", lightCount)

	for i := 0; i < int(lightCount); i++ {
		entry := float32(0)
		err = binary.Read(r, binary.LittleEndian, &entry)
		if err != nil {
			return fmt.Errorf("read entry: %w", err)
		}
		dump.Hex(entry, "%dentry=%0.10f", i, entry)
		e.lights = append(e.lights, entry)
	}

	return nil
}
