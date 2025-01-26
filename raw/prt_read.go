package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Prt is a Particle Render
type Prt struct {
	MetaFileName string
	Version      uint32
	Entries      []*PrtEntry
}

// Identity returns the type of the struct
func (prt *Prt) Identity() string {
	return "prt"
}

// PrtEntry is  ParticleRender entry
type PrtEntry struct {
	ID            uint32
	ID2           uint32
	ParticlePoint string
	//ParticlePointSuffix []byte
	UnknownA1       uint32
	UnknownA2       uint32
	UnknownA3       uint32
	UnknownA4       uint32
	UnknownA5       uint32
	Duration        uint32
	UnknownB        uint32
	UnknownFFFFFFFF int32
	UnknownC        uint32
}

// Read reads a PRT file
func (prt *Prt) Read(r io.ReadSeeker) error {

	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "PTCL" {
		return fmt.Errorf("invalid header %s, wanted EQGS", header)
	}

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

	//fmt.Printf("%s (prt) readd %d entries\n", render.Header.Name, len(render.Entries))
	return nil
}

// SetFileName sets the name of the file
func (prt *Prt) SetFileName(name string) {
	prt.MetaFileName = name
}

// FileName returns the name of the file
func (prt *Prt) FileName() string {
	return prt.MetaFileName
}
