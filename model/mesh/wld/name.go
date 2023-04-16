package wld

import (
	"encoding/binary"
	"fmt"
	"io"
)

func (e *WLD) nameFromHashIndex(r io.ReadSeeker) (string, error) {
	name := ""
	var value int32
	err := binary.Read(r, binary.LittleEndian, &value)
	if err != nil {
		return "", fmt.Errorf("read hash index: %w", err)
	}

	name, ok := e.names[-value]
	if !ok {
		return "", fmt.Errorf("hash 0x%x not found in names (len %d)", -value, len(e.names))
	}
	//dump.Hex(value, "name=(%s)", name)
	return name, nil
}
