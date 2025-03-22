package raw

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Write writes a prt file
func (prt *Prt) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.String("PTCL")
	enc.Uint32(uint32(len(prt.Entries)))
	enc.Uint32(prt.Version)

	for _, entry := range prt.Entries {
		enc.Uint32(entry.ID)
		if prt.Version >= 5 {
			enc.Uint32(entry.ID2)
		}

		enc.StringZero(entry.ParticlePoint)

		data, err := hex.DecodeString(entry.ParticleSuffix)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		enc.Bytes(data)
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

	return nil
}
