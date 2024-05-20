package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/tag"
)

// Pts is a particle point
type Pts struct {
	MetaFileName string `yaml:"file_name"`
	Version      uint32
	Entries      []*PtsEntry `yaml:"entries,omitempty"`
}

// Identity returns the type of the struct
func (pts *Pts) Identity() string {
	return "pts"
}

// PtsEntry is a single entry in a particle point
type PtsEntry struct {
	Name        string  `yaml:"name"`
	BoneName    string  `yaml:"bone_name"`
	Translation Vector3 `yaml:"translation"`
	Rotation    Vector3 `yaml:"rotation"`
	Scale       Vector3 `yaml:"scale"`
	//NameSuffix  []byte  `yaml:"name_suffix,omitempty"`
	//BoneSuffix  []byte  `yaml:"bone_suffix,omitempty"`
}

// Read reads a PTS file
func (pts *Pts) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	header := dec.StringFixed(4)
	if header != "EQPT" {
		return fmt.Errorf("invalid header %s, wanted EQPT", header)
	}

	tag.New()

	particleCount := dec.Uint32()
	pts.Version = dec.Uint32()
	if pts.Version != 1 {
		return fmt.Errorf("invalid version %d, wanted 1", pts.Version)
	}
	tag.Add(0, dec.Pos(), "red", "header")

	for i := 0; i < int(particleCount); i++ {
		entry := &PtsEntry{}
		entry.Name = dec.StringZero()
		_ = dec.Bytes(64 - len(entry.Name) - 1) // entry.NameSuffix
		entry.BoneName = dec.StringZero()
		_ = dec.Bytes(64 - len(entry.BoneName) - 1) // entry.BoneSuffix
		entry.Translation.X = dec.Float32()
		entry.Translation.Y = dec.Float32()
		entry.Translation.Z = dec.Float32()

		entry.Rotation.X = dec.Float32()
		entry.Rotation.Y = dec.Float32()
		entry.Rotation.Z = dec.Float32()

		entry.Scale.X = dec.Float32()
		entry.Scale.Y = dec.Float32()
		entry.Scale.Z = dec.Float32()

		pts.Entries = append(pts.Entries, entry)
		tag.AddRand(tag.LastPos(), dec.Pos(), fmt.Sprintf("%d|%s|%s", i, entry.Name, entry.BoneName))
	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	//log.Debugf("%s (pts) readd %d entries", pts.Header.Name, len(pts.Entries))
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
