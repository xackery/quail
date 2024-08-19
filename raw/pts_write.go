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

		enc.Float32(entry.Translation.X)
		enc.Float32(entry.Translation.Y)
		enc.Float32(entry.Translation.Z)

		enc.Float32(entry.Rotation.X)
		enc.Float32(entry.Rotation.Y)
		enc.Float32(entry.Rotation.Z)

		enc.Float32(entry.Scale.X)
		enc.Float32(entry.Scale.Y)
		enc.Float32(entry.Scale.Z)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
