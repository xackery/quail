package prt

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/tag"
)

// Decode decodes a PRT file
func Decode(render *common.ParticleRender, r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "PTCL" {
		return fmt.Errorf("invalid header %s, wanted EQGS", header)
	}

	tag.New()

	particleCount := dec.Uint32()
	version := dec.Uint32()
	render.Header.Version = int(version)
	if version < 4 {
		return fmt.Errorf("invalid version %d, wanted 4+", version)
	}

	for i := 0; i < int(particleCount); i++ {
		entry := &common.ParticleRenderEntry{}
		entry.ID = dec.Uint32()
		if version >= 5 {
			entry.ID2 = dec.Uint32()
		}

		entry.ParticlePoint = dec.StringZero()
		entry.ParticlePointSuffix = dec.Bytes(64 - len(entry.ParticlePoint) - 1)

		entry.UnknownA1 = dec.Uint32()
		entry.UnknownA2 = dec.Uint32()
		entry.UnknownA3 = dec.Uint32()
		entry.UnknownA4 = dec.Uint32()
		entry.UnknownA5 = dec.Uint32()

		entry.Duration = dec.Uint32()
		entry.UnknownB = dec.Uint32()
		entry.UnknownFFFFFFFF = dec.Int32()
		entry.UnknownC = dec.Uint32()

		render.Entries = append(render.Entries, entry)
	}

	if dec.Error() != nil {
		return fmt.Errorf("decode: %w", dec.Error())
	}

	log.Debugf("%s (prt) decoded %d entries", render.Header.Name, len(render.Entries))
	return nil
}
