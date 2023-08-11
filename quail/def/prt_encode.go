package def

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

// Encode writes a prt file
func (render *ParticleRender) PRTEncode(version uint32, w io.Writer) error {

	tag.New()
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("PTCL")
	enc.Uint32(uint32(len(render.Entries)))
	enc.Uint32(version)

	for _, entry := range render.Entries {
		enc.Uint32(entry.ID)
		if version >= 5 {
			enc.Uint32(entry.ID2)
		}

		enc.StringZero(entry.ParticlePoint)
		enc.Bytes(entry.ParticlePointSuffix)
		enc.Uint32(entry.UnknownA1)
		enc.Uint32(entry.UnknownA2)
		enc.Uint32(entry.UnknownA3)
		enc.Uint32(entry.UnknownA4)
		enc.Uint32(entry.UnknownA5)
		enc.Uint32(entry.Duration)
		enc.Uint32(entry.UnknownB)
		enc.Int32(entry.UnknownFFFFFFFF)
		enc.Uint32(entry.UnknownC)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	log.Debugf("%s prt encoded %d entries", render.Name, len(render.Entries))
	return nil
}
