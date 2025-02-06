package raw

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// LOD is level of detail information
// Typical usaeg is like so:
/*
EQLOD
LOD,OBJ_FIREPIT_STMFT,150
LOD,OBJ_FIREPIT_STMFT_LOD1,250
LOD,OBJ_FIREPIT_STMFT_LOD2,400
LOD,OBJ_FIREPIT_STMFT_LOD3,1000
*/
type Lod struct {
	MetaFileName string
	Entries      []*LodEntry
}

// Identity returns the type of the struct
func (lod *Lod) Identity() string {
	return "lod"
}

type LodEntry struct {
	Category   string
	ObjectName string
	Distance   int
}

func (lod *Lod) Read(r io.ReadSeeker) error {

	// read all of r as a string
	buf := make([]byte, 1024)
	var str string
	for {
		n, err := r.Read(buf)
		if err != nil {
			break
		}
		str += string(buf[:n])
	}

	// split by newline
	lines := strings.Split(str, "\n")
	if strings.TrimSpace(lines[0]) != "EQLOD" {
		return fmt.Errorf("header does not match EQLOD, got %s", lines[0])
	}

	lod.Entries = []*LodEntry{}
	// parse each line
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// split into commas
		records := strings.Split(line, ",")
		if len(records) == 1 {
			continue
		}
		entry := &LodEntry{
			Distance: 0,
		}
		if len(records) > 2 {
			val, err := strconv.Atoi(records[2])
			if err != nil {
				return fmt.Errorf("line %d lod %d is not a number", i, val)
			}
			entry.Distance = val
		}
		if len(records) < 2 {
			return fmt.Errorf("line %d expected at least 2 entries, got %d", i, len(records))
		}

		entry.Category = records[0]
		entry.ObjectName = records[1]

		lod.Entries = append(lod.Entries, entry)
	}

	return nil
}

// SetFileName sets the name of the file
func (lod *Lod) SetFileName(name string) {
	lod.MetaFileName = name
}

// FileName returns the name of the file
func (lod *Lod) FileName() string {
	return lod.MetaFileName
}
