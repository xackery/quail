package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/tag"
)

// Prt is a Particle Render
type Prt struct {
	Version uint32      `yaml:"version"`
	Entries []*PrtEntry `yaml:"entries,omitempty"`
}

// PrtEntry is  ParticleRender entry
type PrtEntry struct {
	ID            uint32 `yaml:"id"` //id is actorsemittersnew.edd
	ID2           uint32 `yaml:"id2"`
	ParticlePoint string `yaml:"particle_point"`
	//ParticlePointSuffix []byte `yaml:"particle_point_suffix,omitempty"`
	UnknownA1       uint32 `yaml:"unknowna1"`
	UnknownA2       uint32 `yaml:"unknowna2"`
	UnknownA3       uint32 `yaml:"unknowna3"`
	UnknownA4       uint32 `yaml:"unknowna4"`
	UnknownA5       uint32 `yaml:"unknowna5"`
	Duration        uint32 `yaml:"duration"`
	UnknownB        uint32 `yaml:"unknownb"`
	UnknownFFFFFFFF int32  `yaml:"unknownffffffff"`
	UnknownC        uint32 `yaml:"unknownc"`
}

// Read reads a PRT file
func (prt *Prt) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "PTCL" {
		return fmt.Errorf("invalid header %s, wanted EQGS", header)
	}

	tag.New()

	particleCount := dec.Uint32()
	prt.Version = dec.Uint32()
	if prt.Version < 4 {
		return fmt.Errorf("invalid version %d, wanted 4+", prt.Version)
	}

	for i := 0; i < int(particleCount); i++ {
		entry := &PrtEntry{}
		entry.ID = dec.Uint32()
		if prt.Version >= 5 {
			entry.ID2 = dec.Uint32()
		}

		entry.ParticlePoint = dec.StringZero()
		_ = dec.Bytes(64 - len(entry.ParticlePoint) - 1) // was entry.ParticlePointSuffix

		entry.UnknownA1 = dec.Uint32()
		entry.UnknownA2 = dec.Uint32()
		entry.UnknownA3 = dec.Uint32()
		entry.UnknownA4 = dec.Uint32()
		entry.UnknownA5 = dec.Uint32()

		entry.Duration = dec.Uint32()
		entry.UnknownB = dec.Uint32()
		entry.UnknownFFFFFFFF = dec.Int32()
		entry.UnknownC = dec.Uint32()

		prt.Entries = append(prt.Entries, entry)
	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	//log.Debugf("%s (prt) readd %d entries", render.Header.Name, len(render.Entries))
	return nil
}
