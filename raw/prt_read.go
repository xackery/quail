package raw

import (
	"encoding/binary"
	"encoding/hex"
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

// PrtEntry is  ParticleRender entry
type PrtEntry struct {
	ID              uint32
	ID2             uint32
	ParticlePoint   string
	ParticleSuffix  string
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

// Identity returns the type of the struct
func (prt *Prt) Identity() string {
	return "prt"
}

func (prt *Prt) String() string {
	out := ""
	out += fmt.Sprintf("metafilename: %s\n", prt.MetaFileName)
	out += fmt.Sprintf("version: %d\n", prt.Version)
	out += fmt.Sprintf("entries: %d\n", len(prt.Entries))
	for i, entry := range prt.Entries {
		out += fmt.Sprintf("  %d: %s %d %d %d %d %d\n", i, entry.ParticlePoint, entry.ID, entry.ID2, entry.Duration, entry.UnknownA1, entry.UnknownA2)

	}
	return out

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
		data := dec.Bytes(64 - len(entry.ParticlePoint) - 1) // was entry.ParticlePointSuffix
		entry.ParticleSuffix = hex.EncodeToString(data)

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

	pos := dec.Pos()
	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if pos != endPos {
		if pos < endPos {
			return fmt.Errorf("%d bytes remaining (%d total)", endPos-pos, endPos)
		}

		return fmt.Errorf("read past end of file")
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
