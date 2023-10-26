package pts

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/tag"
)

// Encode writes a pts file
func Encode(point *common.ParticlePoint, version uint32, w io.Writer) error {
	tag.New()
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("EQPT")
	enc.Uint32(uint32(len(point.Entries)))
	enc.Uint32(version)

	for _, entry := range point.Entries {
		enc.StringZero(entry.Name)
		enc.Bytes(entry.NameSuffix)

		enc.StringZero(entry.BoneName)
		enc.Bytes(entry.BoneSuffix)

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

	//log.Debugf("%s pts encoded %d entries", point.Header.Name, len(point.Entries))
	return nil
}
