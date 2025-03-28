package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Pts is a particle point
type Pts struct {
	MetaFileName string
	Version      uint32
	Entries      []*PtsEntry
}

// Identity returns the type of the struct
func (pts *Pts) Identity() string {
	return "pts"
}

// PtsEntry is a single entry in a particle point
type PtsEntry struct {
	Name        string
	BoneName    string
	Translation [3]float32
	Rotation    [3]float32
	Scale       [3]float32
	//NameSuffix  []byte
	//BoneSuffix  []byte
}

// Read reads a PTS file
func (pts *Pts) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQPT" {
		return fmt.Errorf("invalid header %s, wanted EQPT", header)
	}

	particleCount := dec.Uint32()
	pts.Version = dec.Uint32()
	if pts.Version != 1 {
		return fmt.Errorf("invalid version %d, wanted 1", pts.Version)
	}

	for i := 0; i < int(particleCount); i++ {
		entry := &PtsEntry{}
		entry.Name = dec.StringZero()
		_ = dec.Bytes(64 - len(entry.Name) - 1) // entry.NameSuffix
		entry.BoneName = dec.StringZero()
		_ = dec.Bytes(64 - len(entry.BoneName) - 1) // entry.BoneSuffix
		entry.Translation[0] = dec.Float32()
		entry.Translation[1] = dec.Float32()
		entry.Translation[2] = dec.Float32()

		entry.Rotation[0] = dec.Float32()
		entry.Rotation[1] = dec.Float32()
		entry.Rotation[2] = dec.Float32()

		entry.Scale[0] = dec.Float32()
		entry.Scale[1] = dec.Float32()
		entry.Scale[2] = dec.Float32()

		pts.Entries = append(pts.Entries, entry)
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

	return nil
}

// SetFileName sets the name of the file
func (pts *Pts) SetFileName(name string) {
	pts.MetaFileName = name
}

// FileName returns the name of the file
func (pts *Pts) FileName() string {
	return pts.MetaFileName
}
