package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Encode writes a pts file
func (pts *Pts) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQPT")
	enc.Uint32(uint32(len(pts.Entries)))
	enc.Uint32(pts.Version)

	for _, entry := range pts.Entries {
		enc.StringZero(entry.Name)

		enc.Bytes(make([]byte, 64-len(entry.Name)-1)) //enc.Bytes(entry.NameSuffix)

		enc.StringZero(entry.BoneName)
		enc.Bytes(make([]byte, 64-len(entry.BoneName)-1)) // enc.Bytes(entry.BoneSuffix)

		enc.Float32(entry.Translation[0])
		enc.Float32(entry.Translation[1])
		enc.Float32(entry.Translation[2])

		enc.Float32(entry.Rotation[0])
		enc.Float32(entry.Rotation[1])
		enc.Float32(entry.Rotation[2])

		enc.Float32(entry.Scale[0])
		enc.Float32(entry.Scale[1])
		enc.Float32(entry.Scale[2])
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
