package pts

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

// Decode decodes a PTS file
func Decode(point *common.ParticlePoint, r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQPT" {
		return fmt.Errorf("invalid header %s, wanted EQPT", header)
	}

	tag.New()

	particleCount := dec.Uint32()
	version := dec.Uint32()
	if version != 1 {
		return fmt.Errorf("invalid version %d, wanted 1", version)
	}

	for i := 0; i < int(particleCount); i++ {
		entry := common.ParticlePointEntry{}
		entry.Name = dec.StringZero()
		entry.NameSuffix = dec.Bytes(64 - len(entry.Name) - 1)
		entry.Bone = dec.StringZero()
		entry.BoneSuffix = dec.Bytes(64 - len(entry.Bone) - 1)
		entry.Translation.X = dec.Float32()
		entry.Translation.Y = dec.Float32()
		entry.Translation.Z = dec.Float32()

		entry.Rotation.X = dec.Float32()
		entry.Rotation.Y = dec.Float32()
		entry.Rotation.Z = dec.Float32()

		entry.Scale.X = dec.Float32()
		entry.Scale.Y = dec.Float32()
		entry.Scale.Z = dec.Float32()

		point.Entries = append(point.Entries, entry)
	}

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	log.Debugf("%s (pts) decoded %d entries", point.Name, len(point.Entries))
	return nil
}
